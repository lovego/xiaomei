package db

import (
	"database/sql"
	"github.com/bughou-go/xiaomei/config"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var mysqlConns map[string]*sql.DB

func init() {
	if mysqlConns == nil {
		mysqlConns = make(map[string]*sql.DB)
	}
}

func Mysql(name string) *sql.DB {
	mysql := mysqlConns[name]
	if mysql != nil {
		return mysql
	}
	var err error
	mysql, err = sql.Open(`mysql`, config.Mysql()[name])
	if err != nil {
		panic(err.Error())
	}
	if err := mysql.Ping(); err != nil {
		panic(err.Error())
	}
	mysql.SetConnMaxLifetime(time.Minute * 10)
	mysql.SetMaxIdleConns(5)
	mysql.SetMaxOpenConns(50)
	mysqlConns[name] = mysql
	return mysql
}
