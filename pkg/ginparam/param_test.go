package ginparam

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestContext(rawURL string) *gin.Context {
	u, _ := url.Parse(rawURL)
	req := &http.Request{
		URL:    u,
		Header: make(http.Header),
	}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req
	return c
}

func TestParseParam_Default(t *testing.T) {
	c := newTestContext("/api/v1/users/")
	q := ParseParams(c)

	if q.Page != 1 {
		t.Errorf("default page = %d, want 1", q.Page)
	}
	if q.PageSize != 25 {
		t.Errorf("default page_size = %d, want 25", q.PageSize)
	}
	if len(q.Orders) != 0 {
		t.Errorf("default orders length = %d, want 0", len(q.Orders))
	}
}

func TestParseParam_CustomPage(t *testing.T) {
	c := newTestContext("/api/v1/users/?page=3&page_size=50")
	q := ParseParams(c)

	if q.Page != 3 {
		t.Errorf("page = %d, want 3", q.Page)
	}
	if q.PageSize != 50 {
		t.Errorf("page_size = %d, want 50", q.PageSize)
	}
}

func TestParseParam_LegacyOrder(t *testing.T) {
	c := newTestContext("/api/v1/users/?order_by=created_time&order=asc")
	q := ParseParams(c)

	if len(q.Orders) != 1 {
		t.Fatalf("orders length = %d, want 1", len(q.Orders))
	}
	if q.Orders[0].Field != "created_time" || q.Orders[0].Order != "asc" {
		t.Errorf("order = %+v, want field=created_time, order=asc", q.Orders[0])
	}
}

func TestParseParam_CompactMultiOrder(t *testing.T) {
	c := newTestContext("/api/v1/users/?order_by=created_time:desc,username:asc")
	q := ParseParams(c)

	if len(q.Orders) != 2 {
		t.Fatalf("orders length = %d, want 2", len(q.Orders))
	}
	if q.Orders[0].Field != "created_time" || q.Orders[0].Order != "desc" {
		t.Errorf("order[0] = %+v, want created_time:desc", q.Orders[0])
	}
	if q.Orders[1].Field != "username" || q.Orders[1].Order != "asc" {
		t.Errorf("order[1] = %+v, want username:asc", q.Orders[1])
	}
}

func TestParseParam_JSONOrders(t *testing.T) {
	jsonParam := url.QueryEscape(`[{"field":"id","order":"desc"},{"field":"name","order":"asc"}]`)
	c := newTestContext("/api/v1/users/?orders=" + jsonParam)
	q := ParseParams(c)

	if len(q.Orders) != 2 {
		t.Fatalf("orders length = %d, want 2", len(q.Orders))
	}
	if q.Orders[0].Field != "id" || q.Orders[0].Order != "desc" {
		t.Errorf("order[0] = %+v, want id:desc", q.Orders[0])
	}
	if q.Orders[1].Field != "name" || q.Orders[1].Order != "asc" {
		t.Errorf("order[1] = %+v, want name:asc", q.Orders[1])
	}
}

func TestParseParam_MultiQueryArray(t *testing.T) {
	c := newTestContext("/api/v1/users/?order_by=created_time:desc&order_by=name:asc")
	q := ParseParams(c)

	if len(q.Orders) != 2 {
		t.Fatalf("orders length = %d, want 2", len(q.Orders))
	}
	if q.Orders[0].Field != "created_time" || q.Orders[0].Order != "desc" {
		t.Errorf("order[0] = %+v, want created_time:desc", q.Orders[0])
	}
	if q.Orders[1].Field != "name" || q.Orders[1].Order != "asc" {
		t.Errorf("order[1] = %+v, want name:asc", q.Orders[1])
	}
}
