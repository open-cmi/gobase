package migrate

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var gotemplate string = `package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/gobase/essential/migrate"
)

// ChangeMeInstance migrate
type ChangeMeInstance struct {
}

// Up up migrate
func (mi ChangeMeInstance) Up(db *sqlx.DB) error {
	sqlClause := ` + "`" + `
		CREATE TABLE IF NOT EXISTS template (
			id char(64) NOT NULL PRIMARY KEY,
			name VARCHAR(256) NOT NULL unique DEFAULT ''
		)
	` + "`" + `
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi ChangeMeInstance) Down(db *sqlx.DB) error {
	sqlClause := ` + "`" + `DROP TABLE IF EXISTS template` + "`" + `
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         	 "00000000000000",
		Description: 	 "example",
		Ext:         	 "go",
		Instance:    	 ChangeMeInstance{},
		Ignore:      	 false,
		AlterOpertion:   false,
	})
}

`

var sqlup string = ``
var sqldown string = ``

// GenerateOpt generate opt
type GenerateOpt struct {
	Output string
	Format string
	Name   string
}

var name string = ""

func NewGenerateOpt(name string, format string, output string) *GenerateOpt {
	return &GenerateOpt{
		Name:   name,
		Format: format,
		Output: output,
	}
}

// Run run
func (g *GenerateOpt) Run() error {

	name = g.Name
	if g.Format == "" || g.Format == "go" {
		SetMigrateMode("go")
	} else {
		SetMigrateMode("sql")
	}

	if g.Output == "" {
		rt, _ := os.Getwd()
		g.Output = filepath.Join(rt, "migration")
	}

	SetMigrateDir(g.Output)

	t := time.Now()
	date := fmt.Sprintf("%4d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	if MigrateMode == "go" {
		// 从template模版读取字符串，然后替换日期
		wfile := fmt.Sprintf("%s_%s.go", date, name)
		wfilepath := filepath.Join(MigrateDir, wfile)
		wf, err := os.OpenFile(wfilepath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("create file failed in migrations, please confirm migrations directory is existing")
			return err
		}

		content := gotemplate

		newcontent := strings.Replace(content, "00000000000000", date, -1)
		newcontent = strings.Replace(newcontent, "example", name, -1)
		newcontent = strings.Replace(newcontent, "template", name, -1)
		io.WriteString(wf, newcontent)
	} else {
		// 从template模版读取字符串，然后替换日期
		upfile := fmt.Sprintf("%s_%s.up.sql", date, name)
		wfilepath := filepath.Join(MigrateDir, upfile)
		wf, err := os.OpenFile(wfilepath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Printf("create file failed in migrations\n")
			return err
		}
		io.WriteString(wf, sqlup)

		downfile := fmt.Sprintf("%s_%s.down.sql", date, name)
		wfilepath = filepath.Join(MigrateDir, downfile)
		wf, err = os.OpenFile(wfilepath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Printf("create file failed in migrations\n")
			return err
		}
		io.WriteString(wf, sqldown)
	}
	fmt.Printf("generate file successfully!\n")
	return nil
}
