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
