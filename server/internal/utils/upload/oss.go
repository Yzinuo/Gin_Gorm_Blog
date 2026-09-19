// 文件上传服务：当你上传一个文件（如图片）时，文件存储服务（如阿里云 OSS、AWS S3、腾讯云 COS 等）会将文件存储在其服务器上，并返回一个唯一的 URL，指向该文件。
// 把这个url存储在mysql中可以实现后续的文件访问和管理。
package upload

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	g "gin-blog/internal/global"

	"github.com/gabriel-vasile/mimetype"
)

// Object Storage Service
// 面向接口编程，OSS可选七牛云，阿里云，腾讯云等，方便切换使用接口  适配器模式
type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

func NewOSS() OSS {
	switch g.GetConfig().Upload.OssType {
	case "local":
		return &Local{}
	case "qiniuyun":
		return &Qiniuyun{}
	case "aliyun":
		return &Aliyun{}
	default:
		return &Local{}
	}
}

// UploadReader uploads an in-memory or archived file. Archive imports do not
// have multipart.FileHeader values for each image, so this keeps that workflow
// inside the same storage provider configured for normal uploads.
func UploadReader(filename string, reader io.Reader) (string, string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	expected := map[string]string{
		".avif": "image/avif", ".gif": "image/gif", ".jpeg": "image/jpeg",
		".jpg": "image/jpeg", ".png": "image/png", ".webp": "image/webp",
	}[ext]
	if expected == "" {
		return "", "", fmt.Errorf("unsupported image extension %s", ext)
	}
	tmp, err := os.CreateTemp("", "archive-image-*")
	if err != nil {
		return "", "", err
	}
	name := tmp.Name()
	defer os.Remove(name)
	hash := sha256.New()
	size, err := io.Copy(io.MultiWriter(tmp, hash), io.LimitReader(reader, (10<<20)+1))
	closeErr := tmp.Close()
	if err != nil {
		return "", "", err
	}
	if closeErr != nil {
		return "", "", closeErr
	}
	if size <= 0 || size > 10<<20 {
		return "", "", fmt.Errorf("image must be between 1 B and 10 MiB")
	}
	detected, err := mimetype.DetectFile(name)
	if err != nil || detected.String() != expected {
		return "", "", fmt.Errorf("image extension, MIME and magic bytes do not match")
	}
	digest := hex.EncodeToString(hash.Sum(nil))
	key := fmt.Sprintf("uploads/article/%s/%s/%s/source%s", time.Now().Format("2006/01"), digest[:2], digest, ext)
	store, err := NewObjectStore()
	if err != nil {
		return "", "", err
	}
	f, err := os.Open(name)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	info, err := store.Put(context.Background(), key, f, size, PutOptions{ContentType: expected, CacheControl: ImmutableCacheControl, SHA256: digest})
	if err != nil {
		return "", "", err
	}
	return store.PublicURL(info.Key), info.Key, nil
}
