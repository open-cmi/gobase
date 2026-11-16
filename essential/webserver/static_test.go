package webserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNoRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建临时目录和静态文件结构
	tmpDir := t.TempDir()

	// 1. Alias 测试目录：tmpDir/motagent-web
	motagentDir := filepath.Join(tmpDir, "motagent-web")
	_ = os.MkdirAll(filepath.Join(motagentDir, "assets"), 0755)
	_ = os.WriteFile(filepath.Join(motagentDir, "index.html"), []byte("<html>motagent spa</html>"), 0644)
	_ = os.WriteFile(filepath.Join(motagentDir, "assets", "app.js"), []byte("console.log('app');"), 0644)

	// 2. Root 测试目录：tmpDir/public/mota
	publicDir := filepath.Join(tmpDir, "public")
	motaDir := filepath.Join(publicDir, "mota")
	_ = os.MkdirAll(filepath.Join(motaDir, "static"), 0755)
	_ = os.WriteFile(filepath.Join(motaDir, "index.html"), []byte("<html>mota root spa</html>"), 0644)
	_ = os.WriteFile(filepath.Join(motaDir, "static", "style.css"), []byte("body { color: red; }"), 0644)

	// 配置 static 规则
	origConf := gConf
	defer func() { gConf = origConf }()

	gConf.Static = []StaticRule{
		{
			Base:  "/motagent",
			Alias: motagentDir,
			Index: "index.html",
		},
		{
			Base:  "/mota",
			Root:  publicDir,
			Index: "index.html",
		},
	}

	svc := New()
	_ = svc.Init()

	tests := []struct {
		name         string
		url          string
		expectedCode int
		expectedBody string
		checkJSON    bool
	}{
		{
			name:         "API path returns 404 JSON",
			url:          "/api/v1/users",
			expectedCode: http.StatusNotFound,
			checkJSON:    true,
		},
		{
			name:         "Sub-product API path returns 404 JSON",
			url:          "/motagent/api/v1/status",
			expectedCode: http.StatusNotFound,
			checkJSON:    true,
		},
		{
			name:         "Alias mode: serve exact static file",
			url:          "/motagent/assets/app.js",
			expectedCode: http.StatusOK,
			expectedBody: "console.log('app');",
		},
		{
			name:         "Alias mode: directory fallback to index.html",
			url:          "/motagent/",
			expectedCode: http.StatusOK,
			expectedBody: "<html>motagent spa</html>",
		},
		{
			name:         "Alias mode: SPA sub-route fallback to index.html",
			url:          "/motagent/user/login",
			expectedCode: http.StatusOK,
			expectedBody: "<html>motagent spa</html>",
		},
		{
			name:         "Root mode: serve exact static file",
			url:          "/mota/static/style.css",
			expectedCode: http.StatusOK,
			expectedBody: "body { color: red; }",
		},
		{
			name:         "Root mode: directory fallback to index.html",
			url:          "/mota/",
			expectedCode: http.StatusOK,
			expectedBody: "<html>mota root spa</html>",
		},
		{
			name:         "Root mode: SPA sub-route fallback to index.html",
			url:          "/mota/dashboard",
			expectedCode: http.StatusOK,
			expectedBody: "<html>mota root spa</html>",
		},
		{
			name:         "Unmatched path returns 404",
			url:          "/other/unknown",
			expectedCode: http.StatusNotFound,
			expectedBody: "404 not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			svc.Engine.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("expected code %d, got %d for url %s", tc.expectedCode, w.Code, tc.url)
			}

			if tc.checkJSON {
				var res map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
					t.Fatalf("expected valid JSON body, got %s", w.Body.String())
				}
				if res["ret"] != float64(404) {
					t.Fatalf("expected ret=404, got %v", res["ret"])
				}
			} else if tc.expectedBody != "" {
				if w.Body.String() != tc.expectedBody {
					t.Fatalf("expected body %q, got %q", tc.expectedBody, w.Body.String())
				}
			}
		})
	}
}

func TestNoRouteWithoutStatic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	origConf := gConf
	defer func() { gConf = origConf }()

	// 无 static 配置
	gConf.Static = nil

	svc := New()
	_ = svc.Init()

	req := httptest.NewRequest(http.MethodGet, "/motagent/login", nil)
	w := httptest.NewRecorder()
	svc.Engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	if w.Body.String() != "404 not found" {
		t.Fatalf("expected '404 not found', got %s", w.Body.String())
	}
}
