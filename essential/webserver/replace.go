package webserver

import (
	"net/http"
	"sort"
	"strings"
)

type compiledReplaceRule struct {
	from string
	to   string
}

func applyReplacePrefix(p string, rules []compiledReplaceRule) (string, bool) {
	for _, rule := range rules {
		if p == rule.from {
			if rule.to == "" {
				return "/", true
			}
			return rule.to, true
		}
		if strings.HasPrefix(p, rule.from+"/") {
			rest := strings.TrimPrefix(p, rule.from)
			if rule.to == "" || rule.to == "/" {
				return rest, true
			}
			return rule.to + rest, true
		}
	}
	return p, false
}

// Handler 返回挂载了可选路径前缀替换（replace）的 http.Handler，在 Gin 路由匹配前执行。
func (s *Service) Handler() http.Handler {
	if len(gConf.Replace) == 0 {
		return s.Engine
	}

	var rules []compiledReplaceRule
	for _, r := range gConf.Replace {
		from := strings.TrimRight(r.From, "/")
		if from == "" {
			continue
		}
		to := strings.TrimRight(r.To, "/")
		rules = append(rules, compiledReplaceRule{
			from: from,
			to:   to,
		})
	}

	if len(rules) == 0 {
		return s.Engine
	}

	// 按前缀长度降序排序（最长前缀优先）
	sort.Slice(rules, func(i, j int) bool {
		return len(rules[i].from) > len(rules[j].from)
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if newPath, ok := applyReplacePrefix(r.URL.Path, rules); ok {
			r.URL.Path = newPath
			if r.URL.RawPath != "" {
				if newRaw, okRaw := applyReplacePrefix(r.URL.RawPath, rules); okRaw {
					r.URL.RawPath = newRaw
				}
			}
		}
		s.Engine.ServeHTTP(w, r)
	})
}
