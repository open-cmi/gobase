package sqldb

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/open-cmi/gobase/essential/config"
	"github.com/open-cmi/gobase/pkg/database/postgresdb"
	"github.com/open-cmi/gobase/pkg/database/sqlitedb"

	"github.com/jmoiron/sqlx"
)

// gConfDB sql db
var gConfDB *sqlx.DB

// DBConfig database model
type Config struct {
	Type     string `json:"type"`
	File     string `json:"file,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

var gConf Config

// GetDB get db
func GetDB() *sqlx.DB {
	return gConfDB.Unsafe()
}

func LikePlaceHolder(holderIndex int) string {
	return fmt.Sprintf("'%%' || $%d || '%%'", holderIndex)
}

// Parse db init
func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	if gConf.Type == "postgresql" || gConf.Type == "pg" {
		var dbconf postgresdb.Config
		dbconf.Host = gConf.Host
		dbconf.Port = gConf.Port
		dbconf.User = gConf.User
		dbconf.Password = gConf.Password
		dbconf.Database = gConf.Database

		dbi, err := postgresdb.PostgresqlInit(&dbconf)
		if err != nil {
			return err
		}
		gConfDB = dbi
	} else if gConf.Type == "sqlite3" {
		var dbconf sqlitedb.Config
		dbconf.File = gConf.File
		dbconf.User = gConf.User
		dbconf.Password = gConf.Password
		dbconf.Database = gConf.Database

		dbi, err := sqlitedb.SQLite3Init(&dbconf)
		if err != nil {
			return err
		}
		gConfDB = dbi
	} else {
		return errors.New("db is not supported")
	}

	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	config.RegisterConfig("sqldb", Parse, Save)
}
