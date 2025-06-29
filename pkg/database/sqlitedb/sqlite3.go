package sqlitedb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-cmi/gobase/pkg/eyas"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type Config struct {
	File     string
	User     string
	Password string
	Database string
}

// SQLite3Init init
func SQLite3Init(conf *Config) (db *sqlx.DB, err error) {
	dbfile := conf.File
	if !strings.HasPrefix(dbfile, "/") && !strings.HasPrefix(dbfile, ".") {
		dbfile = filepath.Join(eyas.GetWorkingDir(), "data", conf.File)
	}

	// if filename is absolute path, use file name directly
	_, err = os.Stat(dbfile)
	if err != nil && os.IsNotExist(err) {
		var file *os.File
		// 如果文件不存在，先创建一个
		file, err = os.OpenFile(dbfile, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			return
		}
		file.Close()
	}
	var file string
	if conf.User != "" && conf.Password != "" {
		file = fmt.Sprintf("file:%s?mode=rwc&_auth&_auth_user=%s&_auth_pass=%s", dbfile, conf.User, conf.Password)
	} else {
		file = fmt.Sprintf("file:%s?mode=rwc", dbfile)
	}
	db, err = sqlx.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA synchronous=NORMAL;")
	db.Exec("PRAGMA busy_timeout=5000;")
	return db, nil
}
