package handle

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

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
