package handle

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestAssetSpecIsClosed(t *testing.T) {
	kind, prefix, ok := assetSpec("resume.sticker.08")
	require.True(t, ok)
	require.Equal(t, "image", kind)
	require.Equal(t, "managed/resume/sticker-08", prefix)
	_, _, ok = assetSpec("resume.sticker.09")
	require.False(t, ok)
	_, _, ok = assetSpec("../../etc/passwd")
	require.False(t, ok)
}

func uploadHeaderFromFile(t *testing.T, filename, sourcePath string) *multipart.FileHeader {
	t.Helper()
	source, err := os.Open(sourcePath)
	require.NoError(t, err)
	defer source.Close()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	require.NoError(t, err)
	_, err = io.Copy(part, source)
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	request := httptest.NewRequest("POST", "/", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	require.NoError(t, request.ParseMultipartForm(1<<20))
	t.Cleanup(func() { request.MultipartForm.RemoveAll() })
	return request.MultipartForm.File["file"][0]
}

func TestSpoolUploadChecksExtensionAgainstMagic(t *testing.T) {
	assetPath := filepath.Join("..", "..", "..", "front", "public", "images", "Before-960.webp")
	valid := uploadHeaderFromFile(t, "before.webp", assetPath)
	prepared, err := spoolUpload(valid, 1<<20, imageMIMEByExt)
	require.NoError(t, err)
	require.Equal(t, "image/webp", prepared.mime)
	require.NoError(t, os.Remove(prepared.path))

	spoofed := uploadHeaderFromFile(t, "before.png", assetPath)
	_, err = spoolUpload(spoofed, 1<<20, imageMIMEByExt)
	require.ErrorContains(t, err, "不一致")
}

func TestOptimizedResumeModelContract(t *testing.T) {
	path := filepath.Join("..", "..", "..", "front", "public", "resume", "resume-ready.optimized.glb")
	stat, err := os.Stat(path)
	require.NoError(t, err)
	require.LessOrEqual(t, stat.Size(), int64(maxModelDelivery))
	file := &preparedFile{path: path, size: stat.Size(), mime: "model/gltf-binary"}
	report, err := validateGLB(file)
	require.NoError(t, err)
	require.True(t, report.Valid)
}

func TestResumeModelReportsMissingRequiredNode(t *testing.T) {
	sourcePath := filepath.Join("..", "..", "..", "front", "public", "resume", "resume-ready.optimized.glb")
	content, err := os.ReadFile(sourcePath)
	require.NoError(t, err)
	mutated := bytes.ReplaceAll(content, []byte("eye_L"), []byte("eye_X"))
	require.NotEqual(t, content, mutated)
	path := filepath.Join(t.TempDir(), "missing-eye.glb")
	require.NoError(t, os.WriteFile(path, mutated, 0o600))
	report, err := validateGLB(&preparedFile{path: path, size: int64(len(mutated)), mime: "model/gltf-binary"})
	require.False(t, report.Valid)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "eye_L"))
}

func assetTestRouter(t *testing.T, db *gorm.DB) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(g.CTX_DB, db)
		c.Next()
	})
	return router
}

func assetTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "assets.db")), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.ManagedAsset{}, &model.ManagedAssetVersion{}))
	return db
}

func TestPublicManifestSupportsETag(t *testing.T) {
	previous := g.Conf
	g.Conf = &g.Config{}
	g.Conf.Storage.AssetManifestEnabled = true
	t.Cleanup(func() { g.Conf = previous })
	db := assetTestDB(t)
	asset := model.ManagedAsset{Key: "home.xray", Kind: "image_pair"}
	require.NoError(t, db.Create(&asset).Error)
	manifest := storedManifest{Files: []assetObject{
		{Role: "before.delivery", URL: "https://assets.example/b.webp", Width: 960},
		{Role: "after.delivery", URL: "https://assets.example/a.webp", Width: 960},
	}}
	version := model.ManagedAssetVersion{AssetID: asset.ID, Version: "v1", Status: model.AssetStatusPublished, SourceManifest: marshal(storedManifest{}), DeliveryManifest: marshal(manifest), ValidationReport: marshal(validationReport{Valid: true})}
	require.NoError(t, db.Create(&version).Error)
	require.NoError(t, db.Model(&asset).Update("active_version_id", version.ID).Error)
	router := assetTestRouter(t, db)
	router.GET("/api/front/assets", (&Asset{}).PublicManifest)

	first := httptest.NewRecorder()
	router.ServeHTTP(first, httptest.NewRequest("GET", "/api/front/assets?scope=home", nil))
	require.Equal(t, 200, first.Code)
	etag := first.Header().Get("ETag")
	require.NotEmpty(t, etag)
	require.Equal(t, "no-cache", first.Header().Get("Cache-Control"))

	request := httptest.NewRequest("GET", "/api/front/assets?scope=home", nil)
	request.Header.Set("If-None-Match", etag)
	second := httptest.NewRecorder()
	router.ServeHTTP(second, request)
	require.Equal(t, 304, second.Code)
}

func TestPublishKeepsCurrentAndTwoHistoricalVersions(t *testing.T) {
	db := assetTestDB(t)
	asset := model.ManagedAsset{Key: "resume.sticker.01", Kind: "image"}
	require.NoError(t, db.Create(&asset).Error)
	versions := make([]model.ManagedAssetVersion, 4)
	for index := range versions {
		versions[index] = model.ManagedAssetVersion{Model: model.Model{CreatedAt: time.Now().Add(time.Duration(index) * time.Second)}, AssetID: asset.ID, Version: fmt.Sprintf("v%d", index), Status: model.AssetStatusReady, SourceManifest: marshal(storedManifest{}), DeliveryManifest: marshal(storedManifest{}), ValidationReport: marshal(validationReport{Valid: true})}
		require.NoError(t, db.Create(&versions[index]).Error)
	}
	router := assetTestRouter(t, db)
	router.POST("/api/asset/:key/versions/:id/publish", (&Asset{}).Publish)
	for _, version := range versions {
		response := httptest.NewRecorder()
		path := fmt.Sprintf("/api/asset/%s/versions/%d/publish", asset.Key, version.ID)
		router.ServeHTTP(response, httptest.NewRequest("POST", path, nil))
		require.Equal(t, 200, response.Code)
	}

	var refreshed model.ManagedAsset
	require.NoError(t, db.First(&refreshed, asset.ID).Error)
	require.Equal(t, versions[3].ID, *refreshed.ActiveVersionID)
	var rows []model.ManagedAssetVersion
	require.NoError(t, db.Where("asset_id = ?", asset.ID).Order("id").Find(&rows).Error)
	require.Equal(t, model.AssetStatusPurgePending, rows[0].Status)
	require.Equal(t, model.AssetStatusArchived, rows[1].Status)
	require.Equal(t, model.AssetStatusArchived, rows[2].Status)
	require.Equal(t, model.AssetStatusPublished, rows[3].Status)
}
