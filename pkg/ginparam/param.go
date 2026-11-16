package ginparam

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/gobase/pkg/types"
)

// Param 类型别名，指向 types.QueryParam
type Param = types.QueryParam

// GetUser 从 gin.Context 获取 "user" 字段对应的 map (若存在)
func GetUser(c *gin.Context) map[string]interface{} {
	if c == nil {
		return nil
	}
	u, exist := c.Get("user")
	if exist {
		if user, ok := u.(map[string]interface{}); ok {
			return user
		}
	}
	return nil
}

// ParseParams 别名，兼容复数形式调用
// 从 gin.Context 解析出 types.QueryParam (page, page_size, orders)
func ParseParams(c *gin.Context) types.QueryParam {
	var q types.QueryParam

	// 1. 解析 page (默认第 1 页)
	pagestr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pagestr)
	if err != nil || page <= 0 {
		page = 1
	}
	q.Page = page

	// 2. 解析 page_size (默认 25 条)
	pagesizestr := c.DefaultQuery("page_size", "25")
	pagesize, err := strconv.Atoi(pagesizestr)
	if err != nil || pagesize <= 0 {
		pagesize = 25
	}
	q.PageSize = pagesize

	// 3. 解析 orders
	// 3.1 优先解析 JSON 格式: ?orders=[{"field":"created_time","order":"desc"}]
	if ordersJSON := c.Query("orders"); ordersJSON != "" {
		var items []types.OrderItem
		if err := json.Unmarshal([]byte(ordersJSON), &items); err == nil && len(items) > 0 {
			q.Orders = items
			return q
		}
	}

	// 3.2 解析 order_by 参数 (支持多值、逗号分隔: "created_time:desc,username:asc" 及兼容老参数 order=desc 组合)
	q.Orders = parseOrderItems(c)

	return q
}

func parseOrderItems(c *gin.Context) []types.OrderItem {
	orderBys := c.QueryArray("order_by")
	if len(orderBys) == 0 {
		if ob := c.Query("order_by"); ob != "" {
			orderBys = []string{ob}
		}
	}
	if len(orderBys) == 0 {
		return nil
	}

	fallbackOrder := strings.ToLower(strings.TrimSpace(c.DefaultQuery("order", "desc")))
	if fallbackOrder != "asc" && fallbackOrder != "desc" {
		fallbackOrder = "desc"
	}

	var items []types.OrderItem
	for _, entry := range orderBys {
		parts := strings.Split(entry, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if idx := strings.Index(part, ":"); idx != -1 {
				items = append(items, types.OrderItem{
					Field: strings.TrimSpace(part[:idx]),
					Order: strings.TrimSpace(part[idx+1:]),
				})
			} else {
				items = append(items, types.OrderItem{
					Field: part,
					Order: fallbackOrder,
				})
			}
		}
	}
	return items
}
