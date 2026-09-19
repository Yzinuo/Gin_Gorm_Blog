package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStore struct {
	Root    string
	BaseURL string
}

func cleanObjectKey(key string) (string, error) {
	key = strings.ReplaceAll(key, "\\", "/")
	clean := filepath.ToSlash(filepath.Clean(key))
	if clean == "." || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "/") || strings.ContainsRune(clean, '\x00') {
		return "", fmt.Errorf("invalid object key")
	}
	return clean, nil
}

func (s *LocalStore) Put(ctx context.Context, key string, body io.Reader, size int64, opts PutOptions) (ObjectInfo, error) {
	key, err := cleanObjectKey(key)
	if err != nil {
		return ObjectInfo{}, err
	}
	destination := filepath.Join(s.Root, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return ObjectInfo{}, err
	}
	tmp, err := os.CreateTemp(filepath.Dir(destination), ".upload-*")
	if err != nil {
		return ObjectInfo{}, err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	written, copyErr := io.Copy(tmp, io.LimitReader(&contextReader{ctx: ctx, r: body}, size+1))
	closeErr := tmp.Close()
	if copyErr != nil {
		return ObjectInfo{}, copyErr
	}
	if closeErr != nil {
		return ObjectInfo{}, closeErr
	}
	if written != size {
		return ObjectInfo{}, fmt.Errorf("object size mismatch: expected %d, wrote %d", size, written)
	}
	if err := os.Rename(tmpName, destination); err != nil {
		return ObjectInfo{}, err
	}
	return ObjectInfo{Key: key, Size: size, ContentType: opts.ContentType, CacheControl: opts.CacheControl, SHA256: opts.SHA256}, nil
}

func (s *LocalStore) Head(ctx context.Context, key string) (ObjectInfo, error) {
	if err := ctx.Err(); err != nil {
		return ObjectInfo{}, err
	}
	key, err := cleanObjectKey(key)
	if err != nil {
		return ObjectInfo{}, err
	}
	stat, err := os.Stat(filepath.Join(s.Root, filepath.FromSlash(key)))
	if err != nil {
		return ObjectInfo{}, err
	}
	return ObjectInfo{Key: key, Size: stat.Size()}, nil
}

func (s *LocalStore) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	key, err := cleanObjectKey(key)
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(s.Root, filepath.FromSlash(key)))
}

func (s *LocalStore) PublicURL(key string) string {
	base := strings.TrimRight(s.BaseURL, "/")
	if base != "" && !strings.HasPrefix(base, "/") && !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "/" + base
	}
	return base + "/" + strings.TrimLeft(filepath.ToSlash(key), "/")
}

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (r *contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}
