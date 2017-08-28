package mysql

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lovego/xiaomei/config"
)

var mysqlDBs = struct {
	sync.Mutex
	m map[string]*sql.DB
}{m: make(map[string]*sql.DB)}

func DB(name string) *sql.DB {
	mysqlDBs.Lock()
	defer mysqlDBs.Unlock()
	db := mysqlDBs.m[name]
	if db == nil {
		db = newDB(name)
		mysqlDBs.m[name] = db
	}
	return db
}

func newDB(name string) *sql.DB {
	db, err := sql.Open(`mysql`, config.Get(`mysql`).GetString(name))
	if err != nil {
		log.Panic(err)
	}
	if err := db.Ping(); err != nil {
		log.Panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 10)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(50)
	return db
}
