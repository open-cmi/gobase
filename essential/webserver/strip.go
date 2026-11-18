package webserver

import (
	"net/http"
	"strings"
)

// Handler 返回挂载了可选前缀剥离（strip）的 http.Handler，在 Gin 路由匹配前执行。
func (s *Service) Handler() http.Handler {
	prefix := gConf.Strip
	if prefix == "" {
		return s.Engine
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, prefix) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			if r.URL.RawPath != "" {
				r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, prefix)
				if r.URL.RawPath == "" {
					r.URL.RawPath = "/"
				}
			}
		}
		s.Engine.ServeHTTP(w, r)
	})
}
