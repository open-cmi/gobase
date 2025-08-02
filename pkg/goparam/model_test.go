package goparam

import (
	"testing"
)

func TestGetColumn(t *testing.T) {
	var m struct {
		Data string `db:"db,omitempty"`
	}
	m.Data = "1"

	fields := GetColumn(m, []string{})
	if len(fields) != 1 || fields[0] != "db" {
		t.Errorf("field db parse failed")
		return
	}
}

func TestBuildFinalClause(t *testing.T) {
	var param Param
	param.PageParam.Page = 0
	param.PageParam.PageSize = 25
	param.OrderParam.Order = "asc"
	param.OrderParam.OrderBy = "name"

	clause := BuildFinalClause(&param, []string{})
	if clause != " LIMIT 25 OFFSET 0" {
		t.Errorf("fail1")
		return
	}

	clause = BuildFinalClause(&param, []string{"name"})
	if clause != " ORDER BY name asc LIMIT 25 OFFSET 0" {
		t.Errorf("fail2")
		return
	}
}
