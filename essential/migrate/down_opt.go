package migrate

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// DownOpt down operation
type DownOpt struct {
	DB     *sqlx.DB
	Format string
	Input  string
	Count  int
}

func NewDownOpt(db *sqlx.DB, format string, input string, count int) *DownOpt {
	return &DownOpt{
		DB:     db,
		Format: format,
		Input:  input,
		Count:  count,
	}
}

// Run down operation
func (o *DownOpt) Run() error {

	if o.Format == "" || o.Format == "go" {
		SetMigrateMode("go")
	} else {
		SetMigrateMode("sql")
	}

	if o.Input != "" {
		SetMigrateDir(o.Input)
	}

	db := o.DB
	co := &CurrentOpt{
		DB: db,
	}
	migrations := co.GetMigrationList()
	if len(migrations) == 0 {
		fmt.Printf("no migration to down\n")
		return nil
	}

	if o.Count == 0 {
		o.Count = len(migrations)
	}

	for idx := len(migrations) - 1; idx >= 0 && o.Count > 0; idx-- {
		m := migrations[idx]
		fmt.Printf("start to down migrate: %s %s\n", m.Seq, m.Description)

		var err error
		if m.Ext == "sql" {
			err = ExecSQLMigrate(db, &m, "down")
		} else if m.Ext == "go" {
			err = ExecGoMigrate(db, &m, "down")
		}
		if err == nil {
			dbexec := fmt.Sprintf("delete from migrations where seq='%s'", m.Seq)
			db.Exec(dbexec)
			fmt.Println("successfully!!")
		} else {
			fmt.Printf("migrate down failed: %s\n", err.Error())
		}
		o.Count--
	}

	return nil
}
