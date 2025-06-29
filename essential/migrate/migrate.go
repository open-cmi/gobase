package migrate

import (
	"sort"
	"strings"
)

// Register register
func Register(seq *SeqInfo) {
	GoMigrationList = append(GoMigrationList, *seq)
	sort.SliceStable(GoMigrationList, func(i, j int) bool {
		cmp := strings.Compare(GoMigrationList[i].Seq, GoMigrationList[j].Seq)
		return cmp == -1
	})
}
