package utils

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Db  *sql.DB
	err error
)

func InitDB() (*sql.DB, error) {
	Db, err = sql.Open("mysql", "yangzhe:Q1w2e3r$@tcp(120.27.22.159:3306)/bookstore")
	if err == nil {
		return Db, err
	}
	return nil, err
}
