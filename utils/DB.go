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
	Db, err = sql.Open("mysql", "yangzhe:Q1w2e3r$@tcp(120.27.22.159:3306)/gin_blog")
	if err == nil {
		return Db, err
	}
	return nil, err
}

// B2S -- 查询数据库TIMESTAMP类型的时间字段，得到[]uint8类型，将该类型转换为string类型
func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
