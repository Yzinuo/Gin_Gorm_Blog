package handle

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	chaiwebp "github.com/chai2010/webp"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/image/draw"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils/upload"
)

type Asset struct{}

const (
	maxManagedRequest = 64 << 20
	maxImagePixels    = 32_000_000
	maxModelSource    = 50 << 20
	maxModelDelivery  = 6 << 20
)

type assetObject struct {
	Role        string `json:"role"`
	Key         string `json:"key"`
	URL         string `json:"url,omitempty"`
	ContentType string `json:"content_type"`
	Bytes       int64  `json:"bytes"`
	SHA256      string `json:"sha256"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
}

type storedManifest struct {
	Files []assetObject `json:"files"`
}
type validationReport struct {
	Valid    bool     `json:"valid"`
	Checks   []string `json:"checks"`
	Warnings []string `json:"warnings,omitempty"`
}

type versionResponse struct {
	ID            int              `json:"id"`
	AssetKey      string           `json:"asset_key"`
	Kind          string           `json:"kind"`
	Version       string           `json:"version"`
	Status        string           `json:"status"`
	Source        storedManifest   `json:"source_manifest"`
	Delivery      storedManifest   `json:"delivery_manifest"`
	Validation    validationReport `json:"validation_report"`
	CreatedBy     int              `json:"created_by"`
	CreatedAt     time.Time        `json:"created_at"`
	PublishedAt   *time.Time       `json:"published_at,omitempty"`
	FailureReason string           `json:"failure_reason,omitempty"`
}

type preparedFile struct {
	path, filename, ext, mime, hash string
	size                            int64
	width, height                   int
}

var imageMIMEByExt = map[string]string{
	".gif": "image/gif", ".jpeg": "image/jpeg", ".jpg": "image/jpeg",
	".png": "image/png", ".webp": "image/webp",
}

func assetSpec(key string) (kind, prefix string, ok bool) {
	switch key {
	case "home.xray":
		return "image_pair", "managed/home/xray", true
	case "resume.model":
		return "glb", "managed/resume/model", true
	}
	if strings.HasPrefix(key, "resume.sticker.") && len(key) == len("resume.sticker.01") {
		n, err := strconv.Atoi(strings.TrimPrefix(key, "resume.sticker."))
		if err == nil && n >= 1 && n <= 8 {
			return "image", fmt.Sprintf("managed/resume/sticker-%02d", n), true
		}
	}
	return "", "", false
}

func spoolUpload(header *multipart.FileHeader, maxBytes int64, allowed map[string]string) (*preparedFile, error) {
	if header == nil || header.Size <= 0 || header.Size > maxBytes {
		return nil, fmt.Errorf("文件必须在 1 B 到 %d MiB 之间", maxBytes>>20)
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	expected, ok := allowed[ext]
	if !ok {
		return nil, fmt.Errorf("不支持的文件扩展名 %s", ext)
	}
	source, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer source.Close()
	tmp, err := os.CreateTemp("", "managed-asset-*")
	if err != nil {
		return nil, err
	}
	name := tmp.Name()
	defer func() {
		if err != nil {
			os.Remove(name)
		}
	}()
	hash := sha256.New()
	written, err := io.Copy(io.MultiWriter(tmp, hash), io.LimitReader(source, maxBytes+1))
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return nil, err
	}
	if written != header.Size || written > maxBytes {
		os.Remove(name)
		return nil, fmt.Errorf("文件大小与声明不一致")
	}
	detectedMIME, detectErr := mimetype.DetectFile(name)
	if detectErr != nil {
		os.Remove(name)
		return nil, detectErr
	}
	detected := detectedMIME.String()
	if expected == "model/gltf-binary" {
		f, openErr := os.Open(name)
		if openErr != nil {
			os.Remove(name)
			return nil, openErr
		}
		head := make([]byte, 4)
		_, readErr := io.ReadFull(f, head)
		f.Close()
		if readErr != nil || string(head) != "glTF" {
			os.Remove(name)
			return nil, fmt.Errorf("文件不是有效的 GLB")
		}
	} else if detected != expected {
		os.Remove(name)
		return nil, fmt.Errorf("扩展名、MIME 与文件内容不一致")
	}
	return &preparedFile{path: name, filename: filepath.Base(header.Filename), ext: ext, mime: expected, hash: hex.EncodeToString(hash.Sum(nil)), size: written}, nil
}

func decodeImageFile(file *preparedFile) (image.Image, image.Config, error) {
	f, err := os.Open(file.path)
	if err != nil {
		return nil, image.Config{}, err
	}
	defer f.Close()
	var cfg image.Config
	switch file.mime {
	case "image/png":
		cfg, err = png.DecodeConfig(f)
	case "image/jpeg":
		cfg, err = jpeg.DecodeConfig(f)
	case "image/gif":
		cfg, err = gif.DecodeConfig(f)
	case "image/webp":
		cfg, err = chaiwebp.DecodeConfig(f)
	default:
		err = fmt.Errorf("unsupported image type")
	}
	if err != nil {
		return nil, cfg, fmt.Errorf("图片无法解码: %w", err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || int64(cfg.Width)*int64(cfg.Height) > maxImagePixels {
		return nil, cfg, fmt.Errorf("图片像素总量超过安全限制")
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return nil, cfg, err
	}
	var img image.Image
	switch file.mime {
	case "image/png":
		img, err = png.Decode(f)
	case "image/jpeg":
		img, err = jpeg.Decode(f)
	case "image/gif":
		img, err = gif.Decode(f)
	case "image/webp":
		img, err = chaiwebp.Decode(f)
	}
	file.width, file.height = cfg.Width, cfg.Height
	return img, cfg, err
}

func encodeWebPVariant(img image.Image, width, quality int) (*preparedFile, error) {
	bounds := img.Bounds()
	sourceWidth, sourceHeight := bounds.Dx(), bounds.Dy()
	if width <= 0 || width > sourceWidth {
		width = sourceWidth
	}
	height := max(1, sourceHeight*width/sourceWidth)
	var output image.Image = img
	if width != sourceWidth {
		resized := image.NewNRGBA(image.Rect(0, 0, width, height))
		draw.CatmullRom.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)
		output = resized
	}
	tmp, err := os.CreateTemp("", "asset-delivery-*.webp")
	if err != nil {
		return nil, err
	}
	name := tmp.Name()
	if err = chaiwebp.Encode(tmp, output, &chaiwebp.Options{Lossless: false, Quality: float32(quality), Exact: true}); err != nil {
		tmp.Close()
		os.Remove(name)
		return nil, err
	}
	if err = tmp.Close(); err != nil {
		os.Remove(name)
		return nil, err
	}
	f, err := os.Open(name)
	if err != nil {
		os.Remove(name)
		return nil, err
	}
	h := sha256.New()
	size, err := io.Copy(h, f)
	f.Close()
	if err != nil {
		os.Remove(name)
		return nil, err
	}
	return &preparedFile{path: name, ext: ".webp", mime: "image/webp", hash: hex.EncodeToString(h.Sum(nil)), size: size, width: width, height: height}, nil
}

func uploadPrepared(ctx context.Context, store upload.ObjectStore, file *preparedFile, key, role string) (assetObject, error) {
	f, err := os.Open(file.path)
	if err != nil {
		return assetObject{}, err
	}
	defer f.Close()
	info, err := store.Put(ctx, key, f, file.size, upload.PutOptions{ContentType: file.mime, CacheControl: upload.ImmutableCacheControl, SHA256: file.hash})
	if err != nil {
		return assetObject{}, err
	}
	verified, err := store.Head(ctx, key)
	if err == nil && verified.Size != file.size {
		_ = store.Delete(ctx, key)
		return assetObject{}, fmt.Errorf("上传后的对象大小校验失败")
	}
	return assetObject{Role: role, Key: info.Key, URL: store.PublicURL(info.Key), ContentType: file.mime, Bytes: file.size, SHA256: file.hash, Width: file.width, Height: file.height}, nil
}

func validateGLB(file *preparedFile) (validationReport, error) {
	report := validationReport{Checks: []string{"glTF magic", "glTF v2", "chunk boundaries"}}
	f, err := os.Open(file.path)
	if err != nil {
		return report, err
	}
	defer f.Close()
	header := make([]byte, 12)
	if _, err = io.ReadFull(f, header); err != nil {
		return report, fmt.Errorf("GLB header 不完整")
	}
	if string(header[:4]) != "glTF" || binary.LittleEndian.Uint32(header[4:8]) != 2 || int64(binary.LittleEndian.Uint32(header[8:12])) != file.size {
		return report, fmt.Errorf("GLB header 或声明长度无效")
	}
	var jsonChunk []byte
	position := int64(12)
	for position < file.size {
		chunkHeader := make([]byte, 8)
		if _, err = io.ReadFull(f, chunkHeader); err != nil {
			return report, fmt.Errorf("GLB chunk header 不完整")
		}
		length := int64(binary.LittleEndian.Uint32(chunkHeader[:4]))
		chunkType := binary.LittleEndian.Uint32(chunkHeader[4:])
		position += 8
		if length < 0 || position+length > file.size {
			return report, fmt.Errorf("GLB chunk 越界")
		}
		if chunkType == 0x4E4F534A {
			if length > 8<<20 {
				return report, fmt.Errorf("GLB JSON chunk 过大")
			}
			jsonChunk = make([]byte, length)
			if _, err = io.ReadFull(f, jsonChunk); err != nil {
				return report, err
			}
		} else if _, err = f.Seek(length, io.SeekCurrent); err != nil {
			return report, err
		}
		position += length
	}
	if len(jsonChunk) == 0 {
		return report, fmt.Errorf("GLB 缺少 JSON chunk")
	}
	var document struct {
		Nodes []struct {
			Name string `json:"name"`
		} `json:"nodes"`
		Cameras []struct {
			Name string `json:"name"`
		} `json:"cameras"`
		Animations []struct {
			Name string `json:"name"`
		} `json:"animations"`
	}
	if err = json.Unmarshal(bytes.TrimRight(jsonChunk, " \x00"), &document); err != nil {
		return report, fmt.Errorf("GLB JSON 无效: %w", err)
	}
	names := map[string]bool{}
	for _, node := range document.Nodes {
		names[node.Name] = true
	}
	for _, camera := range document.Cameras {
		names[camera.Name] = true
	}
	animations := map[string]bool{}
	for _, animation := range document.Animations {
		animations[animation.Name] = true
	}
	required := []string{"ResumeCamera", "eye_L", "eye_R", "sticker_01_ucas", "sticker_02_caohua", "sticker_03_zhise", "sticker_04_distributed", "sticker_05_scholarship", "sticker_06_bytedance", "sticker_07_research", "sticker_08_github"}
	for _, name := range required {
		if !names[name] {
			return report, fmt.Errorf("GLB 缺少 %s", name)
		}
	}
	if !animations["CameraAction"] {
		return report, fmt.Errorf("GLB 缺少 CameraAction 动画")
	}
	report.Checks = append(report.Checks, "ResumeCamera", "CameraAction", "eye_L/eye_R", "8 sticker nodes")
	report.Valid = true
	return report, nil
}

func marshal(value any) string { data, _ := json.Marshal(value); return string(data) }
func unmarshalVersion(asset model.ManagedAsset, version model.ManagedAssetVersion) versionResponse {
	response := versionResponse{ID: version.ID, AssetKey: asset.Key, Kind: asset.Kind, Version: version.Version, Status: version.Status, CreatedBy: version.CreatedBy, CreatedAt: version.CreatedAt, PublishedAt: version.PublishedAt, FailureReason: version.FailureReason}
	_ = json.Unmarshal([]byte(version.SourceManifest), &response.Source)
	_ = json.Unmarshal([]byte(version.DeliveryManifest), &response.Delivery)
	_ = json.Unmarshal([]byte(version.ValidationReport), &response.Validation)
	return response
}

func currentUserID(c *gin.Context) int {
	if user, err := CurrentUserAuth(c); err == nil && user != nil {
		return user.ID
	}
	return 0
}

func (*Asset) CreateVersion(c *gin.Context) {
	kind, prefix, ok := assetSpec(c.Param("key"))
	if !ok {
		ReturnError(c, g.ErrRequest, "未知的受管资源")
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxManagedRequest)
	if err := c.Request.ParseMultipartForm(1 << 20); err != nil {
		ReturnError(c, g.ErrFileReceive, "上传请求过大或格式无效")
		return
	}
	if c.Request.MultipartForm != nil {
		defer c.Request.MultipartForm.RemoveAll()
	}
	store, err := upload.NewObjectStore()
	if err != nil {
		ReturnError(c, g.ErrFileUpload, err)
		return
	}
	version := uuid.NewString()
	createdKeys := []string{}
	cleanupObjects := func() {
		for _, key := range createdKeys {
			_ = store.Delete(context.Background(), key)
		}
	}
	var sourceFiles, deliveryFiles []assetObject
	report := validationReport{Valid: true}
	imageAllowed := map[string]string{".png": "image/png", ".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".webp": "image/webp"}
	processImage := func(part, objectName string, maxBytes int64, widths []int, minWidth int) (*preparedFile, error) {
		_, header, err := c.Request.FormFile(part)
		if err != nil {
			return nil, fmt.Errorf("缺少 %s 文件", part)
		}
		file, err := spoolUpload(header, maxBytes, imageAllowed)
		if err != nil {
			return nil, err
		}
		img, cfg, err := decodeImageFile(file)
		if err != nil {
			os.Remove(file.path)
			return nil, err
		}
		if cfg.Width < minWidth {
			os.Remove(file.path)
			return nil, fmt.Errorf("%s 宽度至少为 %d", part, minWidth)
		}
		sourceKey := fmt.Sprintf("%s/%s/source/%s%s", prefix, version, objectName, file.ext)
		object, err := uploadPrepared(c, store, file, sourceKey, part+".source")
		if err != nil {
			os.Remove(file.path)
			return nil, err
		}
		createdKeys = append(createdKeys, sourceKey)
		sourceFiles = append(sourceFiles, object)
		quality := 82
		for _, width := range widths {
			variant, err := encodeWebPVariant(img, width, quality)
			if err != nil {
				os.Remove(file.path)
				return nil, err
			}
			key := fmt.Sprintf("%s/%s/delivery/%s-%d.webp", prefix, version, objectName, variant.width)
			object, err := uploadPrepared(c, store, variant, key, part+".delivery")
			os.Remove(variant.path)
			if err != nil {
				os.Remove(file.path)
				return nil, err
			}
			createdKeys = append(createdKeys, key)
			deliveryFiles = append(deliveryFiles, object)
		}
		return file, nil
	}

	switch kind {
	case "image_pair":
		before, err := processImage("before", "before", 10<<20, []int{960, 1672}, 1280)
		if err != nil {
			cleanupObjects()
			ReturnError(c, g.ErrFileUpload, err)
			return
		}
		defer os.Remove(before.path)
		after, err := processImage("after", "after", 10<<20, []int{960, 1672}, 1280)
		if err != nil {
			cleanupObjects()
			ReturnError(c, g.ErrFileUpload, err)
			return
		}
		defer os.Remove(after.path)
		if before.width != after.width || before.height != after.height {
			cleanupObjects()
			ReturnError(c, g.ErrFileUpload, "Before/After 尺寸必须完全一致")
			return
		}
		report.Checks = []string{"magic/MIME match", "pixel limit", "matching dimensions", "responsive WebP variants"}
	case "image":
		file, err := processImage("file", "image", 5<<20, []int{1024}, 1)
		if err != nil {
			cleanupObjects()
			ReturnError(c, g.ErrFileUpload, err)
			return
		}
		defer os.Remove(file.path)
		if file.width > 4096 || file.height > 4096 {
			cleanupObjects()
			ReturnError(c, g.ErrFileUpload, "贴纸最大为 4096×4096")
			return
		}
		for _, object := range deliveryFiles {
			if object.Bytes > 500<<10 {
				cleanupObjects()
				ReturnError(c, g.ErrFileUpload, "贴纸优化后仍超过 500 KiB")
				return
			}
		}
		report.Checks = []string{"magic/MIME match", "pixel limit", "WebP delivery", "500 KiB hard limit"}
	case "glb":
		_, header, formErr := c.Request.FormFile("file")
		if formErr != nil {
			ReturnError(c, g.ErrFileReceive, "缺少 file 文件")
			return
		}
		source, err := spoolUpload(header, maxModelSource, map[string]string{".glb": "model/gltf-binary"})
		if err != nil {
			ReturnError(c, g.ErrFileUpload, err)
			return
		}
		defer os.Remove(source.path)
		report, err = validateGLB(source)
		if err != nil {
			ReturnError(c, g.ErrFileUpload, err)
			return
		}
		delivery := source
		_, deliveryHeader, deliveryFormErr := c.Request.FormFile("delivery")
		if deliveryFormErr == nil {
			delivery, err = spoolUpload(deliveryHeader, maxModelDelivery, map[string]string{".glb": "model/gltf-binary"})
			if err != nil {
				ReturnError(c, g.ErrFileUpload, fmt.Errorf("delivery: %w", err))
				return
			}
			defer os.Remove(delivery.path)
			if _, err = validateGLB(delivery); err != nil {
				ReturnError(c, g.ErrFileUpload, fmt.Errorf("delivery: %w", err))
				return
			}
		} else if !errors.Is(deliveryFormErr, http.ErrMissingFile) {
			ReturnError(c, g.ErrFileReceive, deliveryFormErr)
			return
		} else if source.size > maxModelDelivery {
			ReturnError(c, g.ErrFileUpload, "原始 GLB 超过 6 MiB 时必须同时上传固定脚本生成的 delivery GLB")
			return
		}
		sourceKey := fmt.Sprintf("%s/%s/source/model.glb", prefix, version)
		sourceObject, err := uploadPrepared(c, store, source, sourceKey, "model.source")
		if err != nil {
			cleanupObjects()
			ReturnError(c, g.ErrFileUpload, err)
			return
		}
		createdKeys = append(createdKeys, sourceKey)
		sourceFiles = append(sourceFiles, sourceObject)
		deliveryKey := fmt.Sprintf("%s/%s/delivery/model.glb", prefix, version)
		deliveryObject, err := uploadPrepared(c, store, delivery, deliveryKey, "model.delivery")
		if err != nil {
			cleanupObjects()
			ReturnError(c, g.ErrFileUpload, err)
			return
		}
		createdKeys = append(createdKeys, deliveryKey)
		deliveryFiles = append(deliveryFiles, deliveryObject)
		report.Checks = append(report.Checks, "source retained", "delivery ≤ 6 MiB")
	}

	db := GetDB(c)
	asset := model.ManagedAsset{Key: c.Param("key"), Kind: kind}
	if err := db.Where("`key` = ?", asset.Key).FirstOrCreate(&asset).Error; err != nil {
		cleanupObjects()
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	row := model.ManagedAssetVersion{AssetID: asset.ID, Version: version, Status: model.AssetStatusReady, SourceManifest: marshal(storedManifest{Files: sourceFiles}), DeliveryManifest: marshal(storedManifest{Files: deliveryFiles}), ValidationReport: marshal(report), CreatedBy: currentUserID(c)}
	if err := db.Create(&row).Error; err != nil {
		cleanupObjects()
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, unmarshalVersion(asset, row))
}

func (*Asset) ListVersions(c *gin.Context) {
	kind, _, ok := assetSpec(c.Param("key"))
	if !ok {
		ReturnError(c, g.ErrRequest, "未知的受管资源")
		return
	}
	var asset model.ManagedAsset
	if err := GetDB(c).Where("`key` = ?", c.Param("key")).First(&asset).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ReturnSuccess(c, []versionResponse{})
		return
	} else if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	asset.Kind = kind
	var rows []model.ManagedAssetVersion
	if err := GetDB(c).Where("asset_id = ? AND status <> ?", asset.ID, model.AssetStatusPurgePending).Order("created_at DESC").Limit(20).Find(&rows).Error; err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	responses := make([]versionResponse, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, unmarshalVersion(asset, row))
	}
	ReturnSuccess(c, responses)
}

func (*Asset) GetVersion(c *gin.Context) {
	_, _, ok := assetSpec(c.Param("key"))
	if !ok {
		ReturnError(c, g.ErrRequest, "未知的受管资源")
		return
	}
	var asset model.ManagedAsset
	if err := GetDB(c).Where("`key` = ?", c.Param("key")).First(&asset).Error; err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	var row model.ManagedAssetVersion
	if err := GetDB(c).Where("id = ? AND asset_id = ?", c.Param("id"), asset.ID).First(&row).Error; err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, unmarshalVersion(asset, row))
}

func (*Asset) Publish(c *gin.Context) {
	_, _, ok := assetSpec(c.Param("key"))
	if !ok {
		ReturnError(c, g.ErrRequest, "未知的受管资源")
		return
	}
	db := GetDB(c)
	err := db.Transaction(func(tx *gorm.DB) error {
		var asset model.ManagedAsset
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("`key` = ?", c.Param("key")).First(&asset).Error; err != nil {
			return err
		}
		var target model.ManagedAssetVersion
		if err := tx.Where("id = ? AND asset_id = ?", c.Param("id"), asset.ID).First(&target).Error; err != nil {
			return err
		}
		if target.Status != model.AssetStatusReady && target.Status != model.AssetStatusArchived && target.Status != model.AssetStatusPublished {
			return fmt.Errorf("只有 ready/archived 版本可以发布")
		}
		now := time.Now()
		if asset.ActiveVersionID != nil && *asset.ActiveVersionID != target.ID {
			if err := tx.Model(&model.ManagedAssetVersion{}).Where("id = ?", *asset.ActiveVersionID).Updates(map[string]any{"status": model.AssetStatusArchived}).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&target).Updates(map[string]any{"status": model.AssetStatusPublished, "published_at": &now}).Error; err != nil {
			return err
		}
		if err := tx.Model(&asset).Update("active_version_id", target.ID).Error; err != nil {
			return err
		}
		var keep []int
		if err := tx.Model(&model.ManagedAssetVersion{}).Where("asset_id = ? AND status IN ?", asset.ID, []string{model.AssetStatusPublished, model.AssetStatusArchived}).Order("published_at DESC, created_at DESC").Limit(3).Pluck("id", &keep).Error; err != nil {
			return err
		}
		query := tx.Model(&model.ManagedAssetVersion{}).Where("asset_id = ? AND status = ?", asset.ID, model.AssetStatusArchived)
		if len(keep) > 0 {
			query = query.Where("id NOT IN ?", keep)
		}
		return query.Update("status", model.AssetStatusPurgePending).Error
	})
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	(&Asset{}).GetVersion(c)
}

func (*Asset) DeleteVersion(c *gin.Context) {
	var asset model.ManagedAsset
	if err := GetDB(c).Where("`key` = ?", c.Param("key")).First(&asset).Error; err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	var row model.ManagedAssetVersion
	if err := GetDB(c).Where("id = ? AND asset_id = ?", c.Param("id"), asset.ID).First(&row).Error; err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	if row.Status != model.AssetStatusDraft && row.Status != model.AssetStatusFailed && row.Status != model.AssetStatusReady {
		ReturnError(c, g.ErrRequest, "已发布或历史版本不能直接删除")
		return
	}
	store, err := upload.NewObjectStore()
	if err != nil {
		ReturnError(c, g.ErrFileUpload, err)
		return
	}
	for _, raw := range []string{row.SourceManifest, row.DeliveryManifest} {
		var manifest storedManifest
		_ = json.Unmarshal([]byte(raw), &manifest)
		for _, object := range manifest.Files {
			if err := store.Delete(c, object.Key); err != nil && !errors.Is(err, os.ErrNotExist) {
				ReturnError(c, g.ErrFileUpload, err)
				return
			}
		}
	}
	if err := GetDB(c).Delete(&row).Error; err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, nil)
}

func publicFile(object assetObject) map[string]any {
	return map[string]any{"src": object.URL, "width": object.Width, "height": object.Height, "bytes": object.Bytes, "sha256": object.SHA256}
}

func (*Asset) PublicManifest(c *gin.Context) {
	if !g.GetConfig().Storage.AssetManifestEnabled {
		c.Header("Cache-Control", "no-cache")
		ReturnSuccess(c, gin.H{"enabled": false, "fallback": true, "revision": "disabled", "assets": gin.H{}})
		return
	}
	scope := c.Query("scope")
	keys := []string{}
	switch scope {
	case "home":
		keys = []string{"home.xray"}
	case "resume":
		keys = []string{"resume.model", "resume.sticker.01", "resume.sticker.02", "resume.sticker.03", "resume.sticker.04", "resume.sticker.05", "resume.sticker.06", "resume.sticker.07", "resume.sticker.08"}
	default:
		ReturnError(c, g.ErrRequest, "scope 必须为 home 或 resume")
		return
	}
	var assets []model.ManagedAsset
	if err := GetDB(c).Preload("ActiveVersion").Where("`key` IN ?", keys).Find(&assets).Error; err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	result := gin.H{}
	revisions := []string{}
	for _, asset := range assets {
		if asset.ActiveVersion == nil {
			continue
		}
		var manifest storedManifest
		if json.Unmarshal([]byte(asset.ActiveVersion.DeliveryManifest), &manifest) != nil {
			continue
		}
		entry := gin.H{"version": asset.ActiveVersion.Version}
		switch asset.Kind {
		case "image_pair":
			for _, role := range []string{"before", "after"} {
				variants := []assetObject{}
				for _, object := range manifest.Files {
					if strings.HasPrefix(object.Role, role+".") {
						variants = append(variants, object)
					}
				}
				sort.Slice(variants, func(i, j int) bool { return variants[i].Width < variants[j].Width })
				if len(variants) > 0 {
					item := publicFile(variants[len(variants)-1])
					srcset := []gin.H{}
					for _, v := range variants {
						srcset = append(srcset, gin.H{"src": v.URL, "width": v.Width})
					}
					item["srcset"] = srcset
					entry[role] = item
				}
			}
		case "glb":
			if len(manifest.Files) > 0 {
				entry["model"] = publicFile(manifest.Files[0])
			}
		case "image":
			if len(manifest.Files) > 0 {
				entry["image"] = publicFile(manifest.Files[0])
			}
		}
		result[asset.Key] = entry
		revisions = append(revisions, asset.Key+":"+asset.ActiveVersion.Version)
	}
	sort.Strings(revisions)
	sum := sha256.Sum256([]byte(strings.Join(revisions, "|")))
	revision := hex.EncodeToString(sum[:8])
	etag := "\"" + revision + "\""
	c.Header("Cache-Control", "no-cache")
	c.Header("ETag", etag)
	if c.GetHeader("If-None-Match") == etag {
		c.Status(http.StatusNotModified)
		return
	}
	ReturnSuccess(c, gin.H{"enabled": true, "fallback": len(result) == 0, "revision": revision, "assets": result})
}
