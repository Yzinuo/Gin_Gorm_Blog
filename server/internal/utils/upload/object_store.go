package upload

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	g "gin-blog/internal/global"
)

const ImmutableCacheControl = "public, max-age=31536000, immutable"

type PutOptions struct {
	ContentType  string
	CacheControl string
	SHA256       string
}

type ObjectInfo struct {
	Key          string `json:"key"`
	Size         int64  `json:"size"`
	ETag         string `json:"etag,omitempty"`
	ContentType  string `json:"content_type,omitempty"`
	CacheControl string `json:"cache_control,omitempty"`
	SHA256       string `json:"sha256,omitempty"`
}

type ObjectStore interface {
	Put(ctx context.Context, key string, body io.Reader, size int64, opts PutOptions) (ObjectInfo, error)
	Head(ctx context.Context, key string) (ObjectInfo, error)
	Delete(ctx context.Context, key string) error
	PublicURL(key string) string
}

func NewObjectStore() (ObjectStore, error) {
	switch strings.ToLower(g.GetConfig().Storage.Provider) {
	case "r2":
		return NewR2Store(context.Background(), g.GetConfig().R2)
	case "", "local":
		return &LocalStore{Root: g.GetConfig().Upload.StorePath, BaseURL: g.GetConfig().Upload.Path}, nil
	default:
		return &LegacyStore{oss: NewOSS()}, nil
	}
}

// LegacyStore keeps existing Aliyun/Qiniu deployments usable while all callers
// move to the provider-neutral API. Head metadata is intentionally limited.
type LegacyStore struct{ oss OSS }

func (s *LegacyStore) Put(ctx context.Context, key string, body io.Reader, size int64, opts PutOptions) (ObjectInfo, error) {
	if err := ctx.Err(); err != nil {
		return ObjectInfo{}, err
	}
	uploader, ok := s.oss.(interface {
		UploadReader(string, io.Reader) (string, string, error)
	})
	if !ok {
		return ObjectInfo{}, fmt.Errorf("legacy provider does not support reader uploads")
	}
	url, storedKey, err := uploader.UploadReader(filepath.Base(key), io.LimitReader(body, size))
	if err != nil {
		return ObjectInfo{}, err
	}
	return ObjectInfo{Key: storedKey, Size: size, ETag: url, ContentType: opts.ContentType, CacheControl: opts.CacheControl, SHA256: opts.SHA256}, nil
}
func (s *LegacyStore) Head(context.Context, string) (ObjectInfo, error) {
	return ObjectInfo{}, fmt.Errorf("legacy provider does not support Head")
}
func (s *LegacyStore) Delete(_ context.Context, key string) error { return s.oss.DeleteFile(key) }
func (s *LegacyStore) PublicURL(key string) string {
	conf := g.GetConfig()
	switch strings.ToLower(conf.Storage.Provider) {
	case "aliyun":
		return fmt.Sprintf("https://%s.%s/%s", conf.Aliyun.Bucket, conf.Aliyun.Endpoint, strings.TrimLeft(key, "/"))
	case "qiniuyun":
		return strings.TrimRight(conf.Qiniu.ImgPath, "/") + "/" + strings.TrimLeft(key, "/")
	default:
		return key
	}
}

// UploadFileCompat preserves POST /api/upload's data:string response while it
// uses the new streaming store.
func UploadFileCompat(ctx context.Context, header *multipart.FileHeader, key, contentType, sha string) (string, string, error) {
	store, err := NewObjectStore()
	if err != nil {
		return "", "", err
	}
	f, err := header.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	info, err := store.Put(ctx, key, f, header.Size, PutOptions{ContentType: contentType, CacheControl: ImmutableCacheControl, SHA256: sha})
	if err != nil {
		return "", "", err
	}
	return store.PublicURL(info.Key), info.Key, nil
}
