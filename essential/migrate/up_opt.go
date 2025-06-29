package migrate

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// UpOpt up operation
type UpOpt struct {
	DB     *sqlx.DB
	Input  string
	Format string
	Count  int
}

func NewUpOpt(db *sqlx.DB, format string, input string, count int) *UpOpt {
	return &UpOpt{
		DB:     db,
		Input:  input,
		Format: format,
		Count:  count,
	}
}

// Run up operation run
func (o *UpOpt) Run() error {

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
		DB: o.DB,
	}
	migrations := co.GetMigrationList()

	lo := &ListOpt{}
	filelist := lo.GetMigrationList()

	if o.Count == 0 {
		o.Count = len(filelist)
	}

	// find start index
	startIndex := 0
	if len(migrations) != 0 {
		var find bool
		latest := migrations[len(migrations)-1]
		for idx, fm := range filelist {
			if strings.Compare(latest.Seq, fm.Seq) < 0 {
				startIndex = idx
				find = true
				break
			}
		}
		if !find {
			fmt.Printf("no migrations\n")
			return nil
		}
	}
	var err error
	for idx := startIndex; idx < len(filelist) && o.Count > 0; idx++ {
		fl := filelist[idx]
		if fl.Ignore {
			continue
		}
		fmt.Printf("start to up migrate: %s %s\n", fl.Seq, fl.Description)

		if fl.Ext == "sql" {
			err = ExecSQLMigrate(db, &fl, "up")
		} else if fl.Ext == "go" {
			err = ExecGoMigrate(db, &fl, "up")
		}

		if err != nil {
			fmt.Printf("migrate failed, error: %s\n", err.Error())
			break
		}
		// 支持回滚
		if !fl.AlterOpertion {
			dbexec := fmt.Sprintf("insert into migrations(seq, description, ext) values('%s','%s','%s')",
				fl.Seq, fl.Description, fl.Ext)
			_, err = db.Exec(dbexec)
			if err != nil {
				fmt.Printf("migrate %s %s failed, %s\n", fl.Seq, fl.Description, err.Error())
				break
			}
		}
		fmt.Println("successfully!!")
		o.Count--
	}
	return err
}
