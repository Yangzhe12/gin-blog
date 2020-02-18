package utils

import (
	"database/sql"
	config "gin-blog/conf"
	"time"

	"github.com/gomodule/redigo/redis"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// Db 是Mysql连接对象
	Db  *sql.DB
	err error
	// RedisPool 是redis连接池对象
	RedisPool *redis.Pool
)

// InitDB 初始化MySQL连接
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

// InitRedisPool 初始化redis连接池
func InitRedisPool() *redis.Pool {
	RedisPool = redis.NewPool(func() (redis.Conn, error) {
		redisConn, err := redis.Dial("tcp", config.GetConfiguration().RedisAddress)
		if err != nil {
			return nil, err
		}
		return redisConn, err
	}, 1)
	RedisPool.IdleTimeout = 240 * time.Second
	RedisPool.MaxActive = 3
	return RedisPool
}
