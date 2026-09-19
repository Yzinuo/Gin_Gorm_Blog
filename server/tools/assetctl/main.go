package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	ginblog "gin-blog/internal"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils/upload"
	"gorm.io/gorm"
)

type reportItem struct {
	Source   string `json:"source"`
	Target   string `json:"target,omitempty"`
	Location string `json:"location"`
	SHA256   string `json:"sha256,omitempty"`
	Bytes    int64  `json:"bytes,omitempty"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
}

type dbReference struct {
	Table, Column string
	ID            int
	Value         string
}

var legacyURL = regexp.MustCompile(`https?://[^\s"'<>)]*(?:aliyuncs\.com|img\.heliar\.top)[^\s"'<>)]*|/public/uploaded/[A-Za-z0-9%._/-]+`)
var mediaExt = map[string]bool{".avif": true, ".gif": true, ".glb": true, ".jpeg": true, ".jpg": true, ".png": true, ".webp": true}

func main() {
	mode := flag.String("mode", "migrate", "migrate or gc")
	configPath := flag.String("config", "config.yml", "server config file")
	apply := flag.Bool("apply", false, "perform writes; default is dry-run")
	backupConfirmed := flag.Bool("backup-confirmed", false, "confirm a database backup exists")
	reportPath := flag.String("report", "", "JSON report path; default stdout")
	roots := flag.String("roots", "../front/public,public/uploaded", "comma-separated repository roots")
	flag.Parse()

	conf := g.ReadConfig(*configPath)
	db := ginblog.InitDatabase(conf)
	var items []reportItem
	var err error
	switch *mode {
	case "migrate":
		if *apply && !*backupConfirmed {
			fatal(fmt.Errorf("--apply requires --backup-confirmed"))
		}
		items, err = migrate(context.Background(), db, strings.Split(*roots, ","), *apply)
	case "gc":
		items, err = gc(context.Background(), db, *apply)
	default:
		err = fmt.Errorf("unknown mode %q", *mode)
	}
	if err != nil {
		fatal(err)
	}
	data, _ := json.MarshalIndent(map[string]any{"generated_at": time.Now().UTC(), "dry_run": !*apply, "items": items}, "", "  ")
	if *reportPath == "" {
		fmt.Println(string(data))
		return
	}
	if err := os.WriteFile(*reportPath, data, 0o600); err != nil {
		fatal(err)
	}
}

func fatal(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

func migrate(ctx context.Context, db *gorm.DB, roots []string, apply bool) ([]reportItem, error) {
	items := scanRepository(roots)
	references, err := scanDatabase(db)
	if err != nil {
		return items, err
	}
	store, err := upload.NewObjectStore()
	if err != nil {
		return items, err
	}
	cache := map[string]string{}
	for _, ref := range references {
		matches := legacyURL.FindAllString(ref.Value, -1)
		updated := ref.Value
		for _, oldURL := range matches {
			item := reportItem{Source: oldURL, Location: fmt.Sprintf("%s[%d].%s", ref.Table, ref.ID, ref.Column), Status: "dry-run"}
			if apply {
				newURL := cache[oldURL]
				if newURL == "" {
					target, size, hash, uploadErr := migrateOne(ctx, store, oldURL)
					if uploadErr != nil {
						item.Status = "failed"
						item.Error = uploadErr.Error()
						items = append(items, item)
						continue
					}
					newURL = target
					cache[oldURL] = target
					item.Bytes = size
					item.SHA256 = hash
				}
				item.Target = newURL
				item.Status = "migrated"
				updated = strings.ReplaceAll(updated, oldURL, newURL)
			}
			items = append(items, item)
		}
		if apply && updated != ref.Value {
			if err := db.Table(ref.Table).Where("id = ?", ref.ID).Update(ref.Column, updated).Error; err != nil {
				return items, err
			}
		}
	}
	return items, nil
}

func scanRepository(roots []string) []reportItem {
	items := []reportItem{}
	for _, root := range roots {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}
		_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				items = append(items, reportItem{Source: path, Location: "repository", Status: "failed", Error: err.Error()})
				return nil
			}
			if entry.IsDir() || !mediaExt[strings.ToLower(filepath.Ext(path))] {
				return nil
			}
			stat, statErr := entry.Info()
			if statErr != nil {
				return nil
			}
			items = append(items, reportItem{Source: path, Location: "repository", Bytes: stat.Size(), Status: "inventory"})
			return nil
		})
	}
	return items
}

func scanDatabase(db *gorm.DB) ([]dbReference, error) {
	definitions := []struct {
		table   string
		columns []string
	}{
		{"article", []string{"img", "content"}}, {"config", []string{"value"}}, {"page", []string{"cover"}}, {"user_info", []string{"avatar"}}, {"resume_profile", []string{"document"}},
	}
	refs := []dbReference{}
	for _, definition := range definitions {
		if !db.Migrator().HasTable(definition.table) {
			continue
		}
		for _, column := range definition.columns {
			if !db.Migrator().HasColumn(definition.table, column) {
				continue
			}
			var rows []struct {
				ID    int
				Value string
			}
			if err := db.Table(definition.table).Select("id, `" + column + "` AS value").Where("`" + column + "` <> ''").Scan(&rows).Error; err != nil {
				return refs, err
			}
			for _, row := range rows {
				if legacyURL.MatchString(row.Value) {
					refs = append(refs, dbReference{Table: definition.table, Column: column, ID: row.ID, Value: row.Value})
				}
			}
		}
	}
	return refs, nil
}

func migrateOne(ctx context.Context, store upload.ObjectStore, oldURL string) (string, int64, string, error) {
	tmp, err := os.CreateTemp("", "asset-migrate-*")
	if err != nil {
		return "", 0, "", err
	}
	name := tmp.Name()
	defer os.Remove(name)
	var reader io.ReadCloser
	if strings.HasPrefix(oldURL, "/public/uploaded/") {
		reader, err = os.Open(filepath.Join("public/uploaded", filepath.Base(oldURL)))
	} else {
		request, requestErr := http.NewRequestWithContext(ctx, http.MethodGet, oldURL, nil)
		if requestErr != nil {
			return "", 0, "", requestErr
		}
		client := http.Client{Timeout: 30 * time.Second}
		response, responseErr := client.Do(request)
		if responseErr != nil {
			return "", 0, "", responseErr
		}
		if response.StatusCode/100 != 2 {
			response.Body.Close()
			return "", 0, "", fmt.Errorf("source returned %s", response.Status)
		}
		reader = response.Body
	}
	if err != nil {
		return "", 0, "", err
	}
	defer reader.Close()
	h := sha256.New()
	size, err := io.Copy(io.MultiWriter(tmp, h), io.LimitReader(reader, (50<<20)+1))
	if err != nil {
		return "", 0, "", err
	}
	if size > 50<<20 {
		return "", 0, "", fmt.Errorf("source exceeds 50 MiB")
	}
	if err = tmp.Close(); err != nil {
		return "", 0, "", err
	}
	hash := hex.EncodeToString(h.Sum(nil))
	ext := strings.ToLower(filepath.Ext(strings.Split(oldURL, "?")[0]))
	if !mediaExt[ext] {
		return "", 0, "", fmt.Errorf("unsupported extension")
	}
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	key := fmt.Sprintf("uploads/migrated/%s/%s/%s%s", time.Now().Format("2006/01"), hash[:2], hash, ext)
	f, err := os.Open(name)
	if err != nil {
		return "", 0, "", err
	}
	defer f.Close()
	if _, err = store.Put(ctx, key, f, size, upload.PutOptions{ContentType: contentType, CacheControl: upload.ImmutableCacheControl, SHA256: hash}); err != nil {
		return "", 0, "", err
	}
	head, err := store.Head(ctx, key)
	if err != nil {
		return "", 0, "", err
	}
	if head.Size != size {
		return "", 0, "", fmt.Errorf("Head size mismatch")
	}
	return store.PublicURL(key), size, hash, nil
}

func gc(ctx context.Context, db *gorm.DB, apply bool) ([]reportItem, error) {
	cutoff := time.Now().Add(-7 * 24 * time.Hour)
	var candidates []model.ManagedAssetVersion
	if err := db.Where("status = ? AND updated_at < ?", model.AssetStatusPurgePending, cutoff).Find(&candidates).Error; err != nil {
		return nil, err
	}
	store, err := upload.NewObjectStore()
	if err != nil {
		return nil, err
	}
	items := []reportItem{}
	for _, version := range candidates {
		var asset model.ManagedAsset
		if err := db.First(&asset, version.AssetID).Error; err != nil {
			return items, err
		}
		if asset.ActiveVersionID != nil && *asset.ActiveVersionID == version.ID {
			continue
		}
		var protected []int
		if err := db.Model(&model.ManagedAssetVersion{}).Where("asset_id = ? AND status IN ?", asset.ID, []string{model.AssetStatusPublished, model.AssetStatusArchived}).Order("published_at DESC, created_at DESC").Limit(3).Pluck("id", &protected).Error; err != nil {
			return items, err
		}
		if contains(protected, version.ID) {
			continue
		}
		keys := manifestKeys(version.SourceManifest, version.DeliveryManifest)
		sort.Strings(keys)
		deleteSucceeded := true
		for _, key := range keys {
			item := reportItem{Source: key, Location: asset.Key + ":" + version.Version, Status: "dry-run"}
			if apply {
				if err := store.Delete(ctx, key); err != nil {
					item.Status = "failed"
					item.Error = err.Error()
					deleteSucceeded = false
				} else {
					item.Status = "deleted"
				}
			}
			items = append(items, item)
		}
		if apply && deleteSucceeded {
			if err := db.Delete(&version).Error; err != nil {
				return items, err
			}
		}
	}
	return items, nil
}

func contains(values []int, target int) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
func manifestKeys(raws ...string) []string {
	keys := []string{}
	for _, raw := range raws {
		var manifest struct {
			Files []struct {
				Key string `json:"key"`
			} `json:"files"`
		}
		_ = json.Unmarshal([]byte(raw), &manifest)
		for _, file := range manifest.Files {
			if file.Key != "" {
				keys = append(keys, file.Key)
			}
		}
	}
	return keys
}
