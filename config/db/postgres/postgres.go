package postgres

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/lovego/bsql"
	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/config/db/dburl"
)

var postgresDBs = struct {
	sync.Mutex
	m map[string]*bsql.DB
}{m: make(map[string]*bsql.DB)}

func DB(name string) *bsql.DB {
	postgresDBs.Lock()
	defer postgresDBs.Unlock()
	db := postgresDBs.m[name]
	if db == nil {
		db = newDB(name)
		postgresDBs.m[name] = db
	}
	return db
}

func newDB(name string) *bsql.DB {
	dbUrl := dburl.Parse(config.Get("postgres").GetString(name))
	db, err := sql.Open("postgres", dbUrl.URL.String())
	if err != nil {
		log.Panic(err)
	}
	if err := db.Ping(); err != nil {
		log.Panic(err)
	}
	db.SetMaxOpenConns(dbUrl.MaxOpen)
	db.SetMaxIdleConns(dbUrl.MaxIdle)
	db.SetConnMaxLifetime(dbUrl.MaxLife)
	return bsql.New(db, 5*time.Second)
}
