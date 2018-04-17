package postgres

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/lovego/bsql"
	"github.com/lovego/xiaomei/config"
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
	db, err := sql.Open(`postgres`, config.Get(`postgres`).GetString(name))
	if err != nil {
		log.Panic(err)
	}
	if err := db.Ping(); err != nil {
		log.Panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 10)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(50)
	return &bsql.DB{db, 5 * time.Second}
}
