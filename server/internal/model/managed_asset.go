package model

import "time"

const (
	AssetStatusDraft        = "draft"
	AssetStatusReady        = "ready"
	AssetStatusPublished    = "published"
	AssetStatusArchived     = "archived"
	AssetStatusFailed       = "failed"
	AssetStatusPurgePending = "purge_pending"
)

type ManagedAsset struct {
	Model
	Key             string               `gorm:"size:120;uniqueIndex;not null" json:"key"`
	Kind            string               `gorm:"size:32;not null" json:"kind"`
	ActiveVersionID *int                 `gorm:"index" json:"active_version_id"`
	ActiveVersion   *ManagedAssetVersion `gorm:"foreignKey:ActiveVersionID" json:"active_version,omitempty"`
}

type ManagedAssetVersion struct {
	Model
	AssetID          int        `gorm:"not null;index;uniqueIndex:idx_asset_version" json:"asset_id"`
	Version          string     `gorm:"size:64;not null;uniqueIndex:idx_asset_version" json:"version"`
	Status           string     `gorm:"size:24;not null;index" json:"status"`
	SourceManifest   string     `gorm:"type:text;not null" json:"-"`
	DeliveryManifest string     `gorm:"type:text;not null" json:"-"`
	ValidationReport string     `gorm:"type:text;not null" json:"-"`
	CreatedBy        int        `gorm:"index" json:"created_by"`
	PublishedAt      *time.Time `json:"published_at"`
	FailureReason    string     `gorm:"type:text" json:"failure_reason,omitempty"`
}
