package model

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestResumePersistenceAndAtomicValidation(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "resume.db")), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&ResumeProfile{}))
	doc, err := GetResume(db)
	require.NoError(t, err)
	require.Len(t, doc.Entries, 8)
	doc.Entries[0].Title = "更新后的研究经历"
	doc.Entries[0].Description = strings.Repeat("长期研究与工程实践。", 80)
	doc.Entries[0].Image = "public/uploaded/new-sticker.png"
	require.NoError(t, SaveResume(db, doc))
	saved, err := GetResume(db)
	require.NoError(t, err)
	require.Equal(t, doc, saved)

	// Invalid multi-entry changes must not partially replace a good profile.
	doc.Entries[0].Title = "不应保存"
	doc.Entries[7].ID = doc.Entries[0].ID
	require.Error(t, SaveResume(db, doc))
	after, err := GetResume(db)
	require.NoError(t, err)
	require.Equal(t, saved, after)
	var count int64
	require.NoError(t, db.Model(&ResumeProfile{}).Count(&count).Error)
	require.EqualValues(t, 1, count)
}

func TestResumeRejectsUnsafeOrIncompleteInput(t *testing.T) {
	for _, image := range []string{"javascript:alert(1)", "data:image/png;base64,AA", "//other.example/image.png", "../secret", "https://user:pass@example.com/a.png"} {
		t.Run(image, func(t *testing.T) {
			doc := DefaultResume()
			doc.Entries[0].Image = image
			require.Error(t, ValidateResume(doc))
		})
	}
	doc := DefaultResume()
	doc.Entries = doc.Entries[:7]
	require.Error(t, ValidateResume(doc))
	doc = DefaultResume()
	doc.Entries[0].Description = strings.Repeat("字", 1201)
	require.Error(t, ValidateResume(doc))
	doc = DefaultResume()
	doc.Entries[0].Title = " "
	require.Error(t, ValidateResume(doc))
	doc = DefaultResume()
	doc.Entries[0].Tags = []string{" "}
	require.Error(t, ValidateResume(doc))
}
