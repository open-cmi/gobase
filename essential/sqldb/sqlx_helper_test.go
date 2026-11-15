package sqldb

import (
	"testing"

	"github.com/open-cmi/gobase/pkg/types"
)

type BaseInfo struct {
	ID        string `db:"id"`
	CreatedAt int64  `db:"created_at"`
}

type UserTestModel struct {
	BaseInfo
	Username string `db:"username,omitempty"`
	Email    string `db:"email"`
	Ignored  string `db:"-"`
	NoTag    string
}

func TestExtractDBColumns(t *testing.T) {
	cols := ExtractDBColumns(UserTestModel{}, "created_at")
	expected := []string{"id", "username", "email"}

	if len(cols) != len(expected) {
		t.Fatalf("ExtractDBColumns length = %d, want %d: %v", len(cols), len(expected), cols)
	}

	for i, c := range cols {
		if c != expected[i] {
			t.Errorf("col[%d] = %s, want %s", i, c, expected[i])
		}
	}
}

func TestClausesDirectString(t *testing.T) {
	cols := []string{"id", "username", "email"}

	// 1. InsertValuesClauses
	valClause := GetInsertValuesClauses(cols)
	wantVal := ":id, :username, :email"
	if valClause != wantVal {
		t.Errorf("GetInsertValuesClauses = %q, want %q", valClause, wantVal)
	}

	// 2. UpdateSetClauses (exclude id)
	updateClause := GetUpdateSetClauses(cols, "id")
	wantUpdate := "username = :username, email = :email"
	if updateClause != wantUpdate {
		t.Errorf("GetUpdateSetClauses = %q, want %q", updateClause, wantUpdate)
	}

	// 3. UpsertSetClauses (exclude id)
	upsertClause := GetUpsertSetClauses(cols, "id")
	wantUpsert := "username = EXCLUDED.username, email = EXCLUDED.email"
	if upsertClause != wantUpsert {
		t.Errorf("GetUpsertSetClauses = %q, want %q", upsertClause, wantUpsert)
	}
}

func TestPageAndOrderClauses(t *testing.T) {
	allowedFields := []string{"created_time", "updated_time", "username"}

	// 1. Page clause
	pClause := BuildPageClause(2, 20)
	wantP := " LIMIT 20 OFFSET 20"
	if pClause != wantP {
		t.Errorf("BuildPageClause = %q, want %q", pClause, wantP)
	}

	// 2. Multi-field Order clause
	items := []types.OrderItem{
		{Field: "created_time", Order: "desc"},
		{Field: "username", Order: "ASC"},
	}
	oClause := BuildOrderClause(items, allowedFields)
	wantO := " ORDER BY created_time desc, username asc"
	if oClause != wantO {
		t.Errorf("BuildOrderClause = %q, want %q", oClause, wantO)
	}

	// 3. Multi-field with invalid injection & duplicate
	itemsWithBad := []types.OrderItem{
		{Field: "created_time", Order: "desc"},
		{Field: "drop table users", Order: "desc"},
		{Field: "created_time", Order: "asc"},               // 重复字段跳过
		{Field: "updated_time", Order: "invalid_direction"}, // 默认 desc
	}
	oClauseFilter := BuildOrderClause(itemsWithBad, allowedFields)
	wantFilter := " ORDER BY created_time desc, updated_time desc"
	if oClauseFilter != wantFilter {
		t.Errorf("BuildOrderClause filter = %q, want %q", oClauseFilter, wantFilter)
	}

	// 4. Order clause (all invalid fields)
	invClause := BuildOrderClause([]types.OrderItem{{Field: "foo", Order: "desc"}}, allowedFields)
	if invClause != "" {
		t.Errorf("BuildOrderClause with invalid field should be empty, got %q", invClause)
	}

	// 5. Final clause with QueryParam
	queryParam := types.QueryParam{
		Page:     1,
		PageSize: 25,
		Orders:   items,
	}
	finalClause := BuildFinalClause(queryParam, allowedFields)
	wantFinal := " ORDER BY created_time desc, username asc LIMIT 25 OFFSET 0"
	if finalClause != wantFinal {
		t.Errorf("BuildFinalClause = %q, want %q", finalClause, wantFinal)
	}

	// 5.1 Final clause by orders
	finalByOrders := BuildFinalClauseByOrders(1, 25, items, allowedFields)
	if finalByOrders != wantFinal {
		t.Errorf("BuildFinalClauseByOrders = %q, want %q", finalByOrders, wantFinal)
	}

	// 6. Single order helper
	singleClause := BuildSingleOrderClause("created_time", "asc", allowedFields)
	if singleClause != " ORDER BY created_time asc" {
		t.Errorf("BuildSingleOrderClause = %q, want %q", singleClause, " ORDER BY created_time asc")
	}

	singleFinal := BuildSingleFinalClause(1, 10, "username", "desc", allowedFields)
	if singleFinal != " ORDER BY username desc LIMIT 10 OFFSET 0" {
		t.Errorf("BuildSingleFinalClause = %q, want %q", singleFinal, " ORDER BY username desc LIMIT 10 OFFSET 0")
	}
}
