package sqldb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/open-cmi/gobase/pkg/types"
)

// JSONStringMap implements sql.Scanner and driver.Valuer for PostgreSQL JSONB map[string]string columns.
type JSONStringMap map[string]string

func (m JSONStringMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "{}", err
	}
	return string(b), nil
}

func (m *JSONStringMap) Scan(src interface{}) error {
	if src == nil {
		*m = nil
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return nil
	}
	if len(data) == 0 {
		*m = nil
		return nil
	}
	return json.Unmarshal(data, m)
}

// JSONMap implements sql.Scanner and driver.Valuer for PostgreSQL JSONB map[string]interface{} columns.
type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "{}", err
	}
	return string(b), nil
}

func (m *JSONMap) Scan(src interface{}) error {
	if src == nil {
		*m = nil
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return nil
	}
	if len(data) == 0 {
		*m = nil
		return nil
	}
	return json.Unmarshal(data, m)
}

// JSONStrings implements sql.Scanner and driver.Valuer for PostgreSQL JSONB []string columns.
type JSONStrings []string

func (s JSONStrings) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return "[]", err
	}
	return string(b), nil
}

func (s *JSONStrings) Scan(src interface{}) error {
	if src == nil {
		*s = nil
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return nil
	}
	if len(data) == 0 {
		*s = nil
		return nil
	}
	return json.Unmarshal(data, s)
}

// ExtractDBColumns 递归解析结构体及内嵌结构体上的 `db:"..."` 标签，提取数据库列名列表
// excludes 支持可变字符串或字符串切片 (如 "id" 或 []string{"id"})
func ExtractDBColumns(model interface{}, excludes ...interface{}) []string {
	excludeSet := make(map[string]bool)
	for _, ex := range excludes {
		switch v := ex.(type) {
		case string:
			if v != "" {
				excludeSet[v] = true
			}
		case []string:
			for _, s := range v {
				if s != "" {
					excludeSet[s] = true
				}
			}
		}
	}

	t := reflect.TypeOf(model)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var cols []string
	var walk func(rt reflect.Type)
	walk = func(rt reflect.Type) {
		if rt.Kind() != reflect.Struct {
			return
		}
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if f.Anonymous {
				ft := f.Type
				for ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}
				walk(ft)
				continue
			}
			tag := f.Tag.Get("db")
			if tag == "" || tag == "-" {
				continue
			}
			// 支持处理 "col_name,omitempty" 格式标签
			if idx := strings.Index(tag, ","); idx != -1 {
				tag = tag[:idx]
			}
			tag = strings.TrimSpace(tag)
			if tag == "" || excludeSet[tag] {
				continue
			}
			cols = append(cols, tag)
		}
	}
	walk(t)
	return cols
}

// GetInsertValuesClauses 直接生成 VALUES (:col1, :col2, ...) 占位符字符串子句
func GetInsertValuesClauses(columns []string) string {
	if len(columns) == 0 {
		return ""
	}
	res := make([]string, len(columns))
	for i, c := range columns {
		res[i] = ":" + c
	}
	return strings.Join(res, ", ")
}

// GetUpdateSetClauses 直接生成 UPDATE SET col1 = :col1, col2 = :col2 字符串子句
func GetUpdateSetClauses(columns []string, excludes ...string) string {
	exSet := make(map[string]bool, len(excludes))
	for _, ex := range excludes {
		exSet[ex] = true
	}
	var res []string
	for _, c := range columns {
		if exSet[c] {
			continue
		}
		res = append(res, fmt.Sprintf("%s = :%s", c, c))
	}
	return strings.Join(res, ", ")
}

// GetUpsertSetClauses 直接生成 ON CONFLICT DO UPDATE SET col1 = EXCLUDED.col1, col2 = EXCLUDED.col2 字符串子句
func GetUpsertSetClauses(columns []string, excludes ...string) string {
	exSet := make(map[string]bool, len(excludes))
	for _, ex := range excludes {
		exSet[ex] = true
	}
	var res []string
	for _, c := range columns {
		if exSet[c] {
			continue
		}
		res = append(res, fmt.Sprintf("%s = EXCLUDED.%s", c, c))
	}
	return strings.Join(res, ", ")
}

// OrderItem 排序项类型别名
type OrderItem = types.OrderItem

// ListOption 查询参数类型别名
type ListOption = types.ListOption

// BuildPageClause 生成分页子句: LIMIT ? OFFSET ?
func BuildPageClause(page, pageSize int) string {
	if pageSize <= 0 {
		return ""
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	return fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
}

// BuildOrderClause 安全拼接多字段排序子句 (防 SQL 注入)
// 例如: items = [{Field: "created_at", Order: "desc"}, {Field: "id", Order: "asc"}]
// 生成: " ORDER BY created_at desc, id asc"
func BuildOrderClause(items []types.OrderItem, allowedFields []string) string {
	if len(items) == 0 || len(allowedFields) == 0 {
		return ""
	}

	allowedMap := make(map[string]bool, len(allowedFields))
	for _, f := range allowedFields {
		allowedMap[f] = true
	}

	seen := make(map[string]bool, len(items))
	var clauses []string
	for _, item := range items {
		field := strings.TrimSpace(item.Field)
		if field == "" || !allowedMap[field] || seen[field] {
			continue
		}
		seen[field] = true
		order := strings.ToLower(strings.TrimSpace(item.Order))
		if order != "asc" && order != "desc" {
			order = "desc"
		}
		clauses = append(clauses, fmt.Sprintf("%s %s", field, order))
	}

	if len(clauses) == 0 {
		return ""
	}
	return " ORDER BY " + strings.Join(clauses, ", ")
}

// BuildFinalClause 一键根据 ListOption 指针生成多字段排序 + 分页子句
func BuildFinalClause(param *types.ListOption, allowedFields []string) string {
	if param == nil {
		return ""
	}
	return BuildOrderClause(param.Orders, allowedFields) + BuildPageClause(param.Page, param.PageSize)
}

// BuildFinalClauseByOrders 根据具体的分页参数和排序项列表生成子句
func BuildFinalClauseByOrders(page, pageSize int, items []types.OrderItem, allowedFields []string) string {
	return BuildOrderClause(items, allowedFields) + BuildPageClause(page, pageSize)
}

// BuildSingleOrderClause 便捷函数：单字段排序
func BuildSingleOrderClause(orderBy, order string, allowedFields []string) string {
	if orderBy == "" {
		return ""
	}
	return BuildOrderClause([]types.OrderItem{{Field: orderBy, Order: order}}, allowedFields)
}

// BuildSingleFinalClause 便捷函数：单字段排序 + 分页
func BuildSingleFinalClause(page, pageSize int, orderBy, order string, allowedFields []string) string {
	return BuildSingleOrderClause(orderBy, order, allowedFields) + BuildPageClause(page, pageSize)
}
