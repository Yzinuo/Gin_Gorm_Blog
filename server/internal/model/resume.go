package model

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:embed resume-defaults.json
var resumeDefaults []byte

// ResumeProfile is separate from the public website-config map. All eight
// slots are persisted in one row so an update cannot leave a partial profile.
type ResumeProfile struct {
	ID       int    `gorm:"primaryKey"`
	Document string `gorm:"type:longtext;not null"`
}

type ResumeEntry struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Image       string   `json:"image"`
}

type ResumeDocument struct {
	Entries []ResumeEntry `json:"entries"`
}

func DefaultResume() ResumeDocument {
	var doc ResumeDocument
	if err := json.Unmarshal(resumeDefaults, &doc); err != nil {
		panic(err)
	}
	return doc
}

func ValidateResume(doc ResumeDocument) error {
	if len(doc.Entries) != 8 {
		return errors.New("请保留全部 8 个贴纸位置")
	}
	seen := map[int]bool{}
	for _, entry := range doc.Entries {
		if entry.ID < 1 || entry.ID > 8 || seen[entry.ID] {
			return errors.New("贴纸位置无效或重复")
		}
		seen[entry.ID] = true
		for _, field := range []struct {
			name, value string
			max         int
		}{
			{"标题", entry.Title, 60}, {"英文分类", entry.Category, 60}, {"介绍", entry.Description, 1200},
		} {
			if strings.TrimSpace(field.value) == "" || utf8.RuneCountInString(field.value) > field.max {
				return fmt.Errorf("第 %d 项%s不能为空，且不能超过 %d 字", entry.ID, field.name, field.max)
			}
		}
		if len(entry.Tags) > 5 {
			return fmt.Errorf("第 %d 项最多设置 5 个标签", entry.ID)
		}
		for _, tag := range entry.Tags {
			if strings.TrimSpace(tag) == "" || utf8.RuneCountInString(tag) > 24 {
				return fmt.Errorf("第 %d 项的标签不能为空或超过 24 字", entry.ID)
			}
		}
		if entry.Image == "" { // Empty means the original embedded sticker.
			continue
		}
		u, err := url.Parse(entry.Image)
		if err != nil || len(entry.Image) > 2048 || strings.ContainsAny(entry.Image, "\\\r\n") {
			return errors.New("贴纸图片地址无效")
		}
		if u.IsAbs() {
			if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil {
				return errors.New("贴纸仅支持 HTTP(S) 图片或本站上传图片")
			}
		} else if u.Host != "" || strings.HasPrefix(entry.Image, "//") || u.Path == "" || strings.Contains(u.Path, "..") {
			return errors.New("本站图片路径无效")
		}
	}
	return nil
}

func GetResume(db *gorm.DB) (ResumeDocument, error) {
	var row ResumeProfile
	err := db.First(&row, 1).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return DefaultResume(), nil
	}
	if err != nil {
		return ResumeDocument{}, err
	}
	var doc ResumeDocument
	if err := json.Unmarshal([]byte(row.Document), &doc); err != nil {
		return doc, err
	}
	return doc, ValidateResume(doc)
}

func SaveResume(db *gorm.DB, doc ResumeDocument) error {
	if err := ValidateResume(doc); err != nil {
		return err
	}
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	return db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, DoUpdates: clause.AssignmentColumns([]string{"document"})}).Create(&ResumeProfile{ID: 1, Document: string(data)}).Error
}
