package db

import (
	"database/sql"
	"sync"
	"time"

	"github.com/bughou-go/xiaomei/config"
	_ "github.com/go-sql-driver/mysql"
)

var mysqlConns = struct {
	sync.RWMutex
	m map[string]*sql.DB
}{m: make(map[string]*sql.DB)}

func Mysql(name string) *sql.DB {
	mysqlConns.RLock()
	mysql := mysqlConns.m[name]
	mysqlConns.RUnlock()
	if mysql != nil {
		return mysql
	}
	var err error
	mysql, err = sql.Open(`mysql`, config.DataSource(`mysql`, name))
	if err != nil {
		panic(err.Error())
	}
	if err := mysql.Ping(); err != nil {
		panic(err.Error())
	}
	mysql.SetConnMaxLifetime(time.Minute * 10)
	mysql.SetMaxIdleConns(5)
	mysql.SetMaxOpenConns(50)
	mysqlConns.Lock()
	mysqlConns.m[name] = mysql
	mysqlConns.Unlock()
	return mysql
}
