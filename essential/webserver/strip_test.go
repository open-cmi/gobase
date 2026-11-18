package webserver

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestStripHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmpDir := t.TempDir()
	motagentDir := filepath.Join(tmpDir, "motagent")
	_ = os.MkdirAll(filepath.Join(motagentDir, "assets"), 0755)
	_ = os.WriteFile(filepath.Join(motagentDir, "assets", "app.js"), []byte("console.log('ok');"), 0644)

	origConf := gConf
	defer func() { gConf = origConf }()

	gConf.Strip = "/motagent"
	gConf.Static = []StaticRule{
		{
			Base:  "/",
			Alias: motagentDir,
			Index: "index.html",
		},
	}

	svc := New()
	_ = svc.Init()

	svc.Engine.GET("/api/overview/v1/", func(c *gin.Context) {
		c.String(http.StatusOK, "overview ok")
	})

	h := svc.Handler()

	// 1. /motagent/api/overview/v1/ 剥离 /motagent 后命中 /api/overview/v1/
	req1 := httptest.NewRequest(http.MethodGet, "/motagent/api/overview/v1/", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK || w1.Body.String() != "overview ok" {
		t.Fatalf("expected 200 'overview ok', got %d %q", w1.Code, w1.Body.String())
	}

	// 2. 不带前缀的直连 API 依然正常命中
	req2 := httptest.NewRequest(http.MethodGet, "/api/overview/v1/", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK || w2.Body.String() != "overview ok" {
		t.Fatalf("expected 200 'overview ok', got %d %q", w2.Code, w2.Body.String())
	}

	// 3. 边界保护：/motagent-other/api/overview/v1/ 不应被误剥离
	req3 := httptest.NewRequest(http.MethodGet, "/motagent-other/api/overview/v1/", nil)
	w3 := httptest.NewRecorder()
	h.ServeHTTP(w3, req3)
	if w3.Code == http.StatusOK && w3.Body.String() == "overview ok" {
		t.Fatalf("should not strip /motagent-other")
	}

	// 4. 静态文件请求剥离 /motagent 后命中 static (base: "/")
	req4 := httptest.NewRequest(http.MethodGet, "/motagent/assets/app.js", nil)
	w4 := httptest.NewRecorder()
	h.ServeHTTP(w4, req4)
	if w4.Code != http.StatusOK || w4.Body.String() != "console.log('ok');" {
		t.Fatalf("expected static js, got %d %q", w4.Code, w4.Body.String())
	}
}
