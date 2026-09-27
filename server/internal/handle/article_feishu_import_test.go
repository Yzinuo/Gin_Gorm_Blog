package handle

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestImportFeishuArchiveCreatesDraftWithChosenMetadata(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "feishu.db")), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Article{}, &model.ArticleTag{}, &model.Category{}, &model.Tag{}, &model.Config{}))
	router := feishuTestRouter(db, &model.UserAuth{Model: model.Model{ID: 42}, IsSuper: true})

	response := postFeishuArchive(t, router, map[string]string{
		"title": "自定义标题", "category_name": "工程实践", "tag_names": `["飞书","笔记"]`,
	}, makeZIP(t, map[string]string{"note.md": "# 导入正文"}))
	require.EqualValues(t, 0, response["code"])
	var article model.Article
	require.NoError(t, db.Preload("Category").Preload("Tags").First(&article).Error)
	require.Equal(t, "自定义标题", article.Title)
	require.Equal(t, "# 导入正文", article.Content)
	require.Equal(t, "工程实践", article.Category.Name)
	require.Equal(t, model.STATUS_DRAFT, article.Status)
	require.Equal(t, 42, article.UserId)
	require.ElementsMatch(t, []string{"飞书", "笔记"}, []string{article.Tags[0].Name, article.Tags[1].Name})
}

func TestImportFeishuArchiveRequiresExistingImportPermission(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "feishu-permission.db")), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Article{}, &model.ArticleTag{}, &model.Category{}, &model.Tag{}, &model.Config{}, &model.Resource{}, &model.RoleResource{}))
	auth := &model.UserAuth{Model: model.Model{ID: 7}, Roles: []*model.Role{{Model: model.Model{ID: 3}}}}
	router := feishuTestRouter(db, auth)
	fields := map[string]string{"title": "权限检查", "category_name": "工程", "tag_names": `["笔记"]`}
	archive := makeZIP(t, map[string]string{"note.md": "正文"})
	require.NotEqualValues(t, 0, postFeishuArchive(t, router, fields, archive)["code"])
	resource := model.Resource{Url: "/article/import", Method: http.MethodPost, Name: "导入文章"}
	require.NoError(t, db.Create(&resource).Error)
	require.NoError(t, db.Create(&model.RoleResource{RoleId: 3, ResourceId: resource.ID}).Error)
	require.EqualValues(t, 0, postFeishuArchive(t, router, fields, archive)["code"])
}

func feishuTestRouter(db *gorm.DB, auth *model.UserAuth) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("feishu-test", cookie.NewStore([]byte("test-secret"))))
	router.Use(func(c *gin.Context) {
		c.Set(g.CTX_DB, db)
		c.Set(g.CTX_USER_AUTH, auth)
		c.Next()
	})
	router.POST("/import/feishu", (&Article{}).ImportFeishuArchive)
	return router
}

func postFeishuArchive(t *testing.T, router *gin.Engine, fields map[string]string, archive []byte) map[string]any {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range fields {
		require.NoError(t, writer.WriteField(key, value))
	}
	file, err := writer.CreateFormFile("file", "feishu.zip")
	require.NoError(t, err)
	_, err = file.Write(archive)
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	request := httptest.NewRequest(http.MethodPost, "/import/feishu", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	var result map[string]any
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
	return result
}

func TestReadFeishuArchiveFindsMarkdownAndAssets(t *testing.T) {
	archive := makeZIP(t, map[string]string{
		"Feishu note.md":     "![diagram](assets/diagram.png)",
		"assets/diagram.png": "image-data",
	})

	content, markdownPath, files, err := readFeishuArchive(archive)

	require.NoError(t, err)
	require.Equal(t, "Feishu note.md", markdownPath)
	require.Equal(t, "![diagram](assets/diagram.png)", content)
	require.Contains(t, files, "assets/diagram.png")
}

func TestReadFeishuArchiveRejectsMultipleMarkdownFiles(t *testing.T) {
	archive := makeZIP(t, map[string]string{
		"first.md":  "first",
		"second.md": "second",
	})

	_, _, _, err := readFeishuArchive(archive)

	require.ErrorContains(t, err, "只能包含一个 Markdown")
}

func TestArchiveImageReference(t *testing.T) {
	imagePath, isLocal := archiveImageReference("exports", "assets/diagram%20one.png")
	require.True(t, isLocal)
	require.Equal(t, "exports/assets/diagram one.png", imagePath)

	_, isLocal = archiveImageReference("exports", "https://example.com/diagram.png")
	require.False(t, isLocal)

	_, isLocal = archiveImageReference("exports", "../../outside.png")
	require.False(t, isLocal)
}

func makeZIP(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var output bytes.Buffer
	writer := zip.NewWriter(&output)
	for filePath, content := range files {
		file, err := writer.Create(filePath)
		require.NoError(t, err)
		_, err = file.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())
	return output.Bytes()
}
