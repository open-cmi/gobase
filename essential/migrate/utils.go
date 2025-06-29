package migrate

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

// MigrateMode mode
var MigrateMode string = "go"

// MigrateDir migreate dir
var MigrateDir string = ""

// SetMigrateMode set migrate mode
func SetMigrateMode(mode string) {
	MigrateMode = mode
}

// SetMigrateDir set migrate directory
func SetMigrateDir(dir string) {
	MigrateDir = dir
}

// ExecSQLMigrate exec sql mod
func ExecSQLMigrate(db *sqlx.DB, si *SeqInfo, updown string) (err error) {
	sqlfile := si.Seq + "_" + si.Description + "." + updown + "." + si.Ext
	sqlfilepath := filepath.Join(MigrateDir, sqlfile)

	_, err = os.Stat(sqlfilepath)
	if err != nil {
		return fmt.Errorf("migrate file %s is not existing", sqlfilepath)
	}

	// exec file content
	f, err := os.Open(sqlfilepath)
	if err != nil {
		return fmt.Errorf("open %s failed", sqlfilepath)
	}

	sqlContent, err := io.ReadAll(f)
	if err != nil {
		errmsg := fmt.Sprintf("read %s failed\n", sqlfilepath)
		return errors.New(errmsg)
	}

	arr := strings.SplitAfter(string(sqlContent), ";")
	for _, sentence := range arr {
		if strings.Trim(sentence, "") == "" {
			continue
		}
		_, err = db.Exec(sentence)
		if err != nil {
			errmsg := fmt.Sprintf("migrate failed %s\n", err.Error())
			return errors.New(errmsg)
		}
	}
	return
}

// ExecGoMigrate exec go migrate
func ExecGoMigrate(db *sqlx.DB, si *SeqInfo, updown string) (err error) {
	instance := si.Instance
	ref := reflect.ValueOf(instance)
	var fun reflect.Value
	if updown == "up" {
		fun = ref.MethodByName("Up")
	} else if updown == "down" {
		fun = ref.MethodByName("Down")
	}
	var params []reflect.Value = []reflect.Value{reflect.ValueOf(db)}
	retlist := fun.Call(params)
	if retlist[0].Interface() != nil {
		return retlist[0].Interface().(error)
	}
	return nil
}
