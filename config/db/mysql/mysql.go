package mysql

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lovego/xiaomei/config"
)

var mysqlConns = struct {
	sync.Mutex
	m map[string]*sql.DB
}{m: make(map[string]*sql.DB)}

func DB(name string) *sql.DB {
	mysqlConns.Lock()
	defer mysqlConns.Unlock()
	mysql := mysqlConns.m[name]
	if mysql != nil {
		return mysql
	}
	var err error
	mysql, err = sql.Open(`mysql`, config.Get(`mysql`).GetString(name))
	if err != nil {
		panic(err.Error())
	}
	if err := mysql.Ping(); err != nil {
		panic(err.Error())
	}
	mysql.SetConnMaxLifetime(time.Minute * 10)
	mysql.SetMaxIdleConns(5)
	mysql.SetMaxOpenConns(50)
	mysqlConns.m[name] = mysql
	return mysql
}
