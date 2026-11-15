package types

// OrderItem 排序项定义
type OrderItem struct {
	Field string `json:"field"` // 排序字段名
	Order string `json:"order"` // 排序方向: "asc" 或 "desc" (不区分大小写，默认 "desc")
}

type QueryParam struct {
	Page     int
	PageSize int
	Orders   []OrderItem
}
