package migrate

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// CurrentOpt get current migration
type CurrentOpt struct {
	DB *sqlx.DB
}

func NewCurrentOpt(db *sqlx.DB) *CurrentOpt {
	return &CurrentOpt{
		DB: db,
	}
}

// GetInstanceFromSeq get instance from seq
func GetInstanceFromSeq(seq string) interface{} {
	for _, ins := range GoMigrationList {
		if ins.Seq == seq {
			return ins.Instance
		}
	}
	return nil
}

// GetMigrationList get migration list
func (o *CurrentOpt) GetMigrationList() (migrations []SeqInfo) {
	db := o.DB
	dbquery := `select * from migrations order by seq asc`
	r, err := db.Query(dbquery)
	if err != nil {
		// select error , init table
		var init string = `CREATE TABLE IF NOT EXISTS migrations (
							seq varchar(14) UNIQUE NOT NULL,
							description varchar(100) NOT NULL default '',
							ext varchar(100) NOT NULL default 'sql'
							)`
		_, err = db.Exec(init)
		if err != nil {
			fmt.Println("init migrations table failed")
		}
		return
	}

	for r.Next() {
		var row SeqInfo
		err = r.Scan(&row.Seq, &row.Description, &row.Ext)
		if err != nil {
			break
		}
		if row.Ext == "go" {
			row.Instance = GetInstanceFromSeq(row.Seq)
			if row.Instance == nil {
				fmt.Printf("\033[1;31;40mWarning\033[0m seq %s's Ext is go, but it's instance is not found\n", row.Seq)
			}
		}
		migrations = append(migrations, row)
	}
	return migrations
}

// Run run
func (o *CurrentOpt) Run() error {

	migrations := o.GetMigrationList()
	if len(migrations) == 0 {
		fmt.Printf("no migrations\n")
		return nil
	}
	for _, m := range migrations {
		fmt.Printf("%s %s %s\n", m.Seq, m.Description, m.Ext)
	}
	return nil
}
