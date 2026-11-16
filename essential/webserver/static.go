package webserver

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/gobase/pkg/eyas"
)

func (s *Service) initNoRoute() {
	s.Engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 2. 若未配置静态规则，作为纯 API 服务直接返回 404，不进行任何磁盘 IO
		if len(gConf.Static) == 0 {
			c.String(http.StatusNotFound, "404 not found")
			return
		}

		// 1. API 请求保护：若包含 /api/，则认定为后端 API 请求，直接返回 JSON 404，不进入静态文件处理
		if strings.HasPrefix(path, "/api/") || strings.Contains(path, "/api/") {
			c.String(http.StatusNotFound, "404 api not found")
			return
		}

		// 3. 按照 Base 长度降序排序（最长前缀优先匹配）
		rules := make([]StaticRule, len(gConf.Static))
		copy(rules, gConf.Static)
		sort.Slice(rules, func(i, j int) bool {
			return len(rules[i].Base) > len(rules[j].Base)
		})

		workDir := eyas.GetWorkingDir()

		for _, rule := range rules {
			base := strings.TrimRight(rule.Base, "/")
			if base == "" {
				base = "/"
			}

			matched := false
			var relPath string
			if base == "/" {
				matched = true
				relPath = strings.TrimPrefix(path, "/")
			} else if path == base {
				matched = true
				relPath = ""
			} else if strings.HasPrefix(path, base+"/") {
				matched = true
				relPath = strings.TrimPrefix(path, base+"/")
			}

			if !matched {
				continue
			}

			indexFile := rule.Index
			if indexFile == "" {
				indexFile = "index.html"
			}

			var targetDir string
			var fpath string
			var fallbackFile string

			if rule.Alias != "" {
				// Alias 模式：用 Alias 目录替换 Base 路径
				if filepath.IsAbs(rule.Alias) {
					targetDir = filepath.Clean(rule.Alias)
				} else {
					targetDir = filepath.Clean(filepath.Join(workDir, rule.Alias))
				}
				fpath = filepath.Clean(filepath.Join(targetDir, relPath))
				fallbackFile = filepath.Clean(filepath.Join(targetDir, indexFile))
			} else if rule.Root != "" {
				// Root 模式：完整 URL 路径追加到 Root 目录之后
				if filepath.IsAbs(rule.Root) {
					targetDir = filepath.Clean(rule.Root)
				} else {
					targetDir = filepath.Clean(filepath.Join(workDir, rule.Root))
				}
				fpath = filepath.Clean(filepath.Join(targetDir, path))
				fallbackFile = filepath.Clean(filepath.Join(targetDir, base, indexFile))
			} else {
				continue
			}

			// 防止路径穿越攻击 (../)
			if !strings.HasPrefix(fpath, targetDir) {
				c.String(http.StatusForbidden, "forbidden")
				return
			}

			// try_files 1: 物理文件存在且不是目录
			if fi, err := os.Stat(fpath); err == nil && !fi.IsDir() {
				c.File(fpath)
				return
			}

			// try_files 2: 访问目录时尝试返回目录下的 index 文件
			indexUnderDir := filepath.Clean(filepath.Join(fpath, indexFile))
			if strings.HasPrefix(indexUnderDir, targetDir) {
				if fi, err := os.Stat(indexUnderDir); err == nil && !fi.IsDir() {
					c.File(indexUnderDir)
					return
				}
			}

			// try_files 3: SPA 单页应用前端路由回退 (fallback 到 indexFile)
			if strings.HasPrefix(fallbackFile, targetDir) {
				if fi, err := os.Stat(fallbackFile); err == nil && !fi.IsDir() {
					c.File(fallbackFile)
					return
				}
			}

			c.String(http.StatusNotFound, "page not found")
			return
		}

		// 4. 未命中任何 static 规则
		c.String(http.StatusNotFound, "404 not found")
	})
}
