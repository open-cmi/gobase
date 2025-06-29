package postgresdb

import (
	"fmt"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

var maxConnections int = 50

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// PostgresqlInit init
func PostgresqlInit(conf *Config) (db *sqlx.DB, err error) {
	host := conf.Host
	port := conf.Port
	user := conf.User
	password := conf.Password
	database := conf.Database

	dbstr := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable", user, password, host, port)
	db, err = sqlx.Open("postgres", dbstr)
	if err != nil {
		return nil, err
	}

	dbstr = fmt.Sprintf("select datname from pg_database where datname='%s'", database)
	row := db.QueryRow(dbstr)
	var dat string
	err = row.Scan(&dat)
	if err != nil {
		// database is is not existing, create
		createdb := fmt.Sprintf("create database %s", database)
		_, err = db.Exec(createdb)
		if err != nil {
			return nil, err
		}
	}
	db.Close()

	dbstr = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, database)
	db, err = sqlx.Open("postgres", dbstr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConnections)
	return db, nil
}
