package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var mysql *sql.DB

func Mysql() *sql.DB {
	if mysql != nil {
		return mysql
	}
	var err error
	mysql, err = sql.Open(`mysql`, Data.Mysql)
	if err != nil {
		panic(err.Error())
	}
	if err := mysql.Ping(); err != nil {
		panic(err.Error())
	}
	return mysql
}
