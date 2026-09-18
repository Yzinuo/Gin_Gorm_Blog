package handle

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestResumeEndpointsAndLegacyAbout(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "resume.db")), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.ResumeProfile{}, &model.Config{}))
	require.NoError(t, model.FindOrCreateConfig(db, g.CONFIG_ABOUT, "原来的 Markdown"))
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(g.CTX_DB, db)
		if c.GetHeader("X-Test-Authenticated") == "yes" {
			c.Set(g.CTX_USER_AUTH, &model.UserAuth{})
		}
	})
	api := &BlogInfo{}
	r.GET("/about", api.GetAbout)
	r.PUT("/about", api.UpdateAbout)
	call := func(method, target string, doc any, authenticated bool) map[string]any {
		body, err := json.Marshal(doc)
		require.NoError(t, err)
		req := httptest.NewRequest(method, target, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if authenticated {
			req.Header.Set("X-Test-Authenticated", "yes")
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var response map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		return response
	}
	require.Equal(t, "原来的 Markdown", call(http.MethodGet, "/about", nil, false)["data"])
	doc := model.DefaultResume()
	doc.Entries[0].Title = "后台更新"
	require.NotEqual(t, float64(0), call(http.MethodPut, "/about?view=resume", doc, false)["code"])
	require.Equal(t, float64(0), call(http.MethodPut, "/about?view=resume", doc, true)["code"])
	result := call(http.MethodGet, "/about?view=resume", nil, false)
	require.Equal(t, float64(0), result["code"])
	entries := result["data"].(map[string]any)["entries"].([]any)
	require.Equal(t, "后台更新", entries[0].(map[string]any)["title"])
	require.Equal(t, "原来的 Markdown", call(http.MethodGet, "/about", nil, false)["data"])
	doc.Entries[0].Title = ""
	require.NotEqual(t, float64(0), call(http.MethodPut, "/about?view=resume", doc, true)["code"])
}
