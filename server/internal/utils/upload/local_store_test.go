package upload

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalStoreContract(t *testing.T) {
	store := &LocalStore{Root: t.TempDir(), BaseURL: "/public/uploaded"}
	body := []byte("immutable asset")
	info, err := store.Put(context.Background(), "managed/test/v1/file.txt", bytes.NewReader(body), int64(len(body)), PutOptions{ContentType: "text/plain", CacheControl: ImmutableCacheControl, SHA256: "abc"})
	require.NoError(t, err)
	require.Equal(t, "managed/test/v1/file.txt", info.Key)
	require.Equal(t, "/public/uploaded/managed/test/v1/file.txt", store.PublicURL(info.Key))

	head, err := store.Head(context.Background(), info.Key)
	require.NoError(t, err)
	require.Equal(t, int64(len(body)), head.Size)
	require.NoError(t, store.Delete(context.Background(), info.Key))
	_, err = store.Head(context.Background(), info.Key)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestLocalStoreRejectsTraversal(t *testing.T) {
	store := &LocalStore{Root: t.TempDir(), BaseURL: "/files"}
	_, err := store.Put(context.Background(), "../secret", bytes.NewReader([]byte("x")), 1, PutOptions{})
	require.Error(t, err)
}
