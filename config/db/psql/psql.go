package psql

import (
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/lovego/xiaomei/config"
)

var postgresDBs = struct {
	sync.Mutex
	m map[string]*sql.DB
}{m: make(map[string]*sql.DB)}

func DB(name string) *sql.DB {
	postgresDBs.Lock()
	defer postgresDBs.Unlock()
	db := postgresDBs.m[name]
	if db == nil {
		db = newDB(name)
		postgresDBs.m[name] = db
	}
	return db
}

func newDB(name string) *sql.DB {
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
	return db
}

func RunInTransaction(db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func Q(q string) string {
	return `'` + strings.Replace(q, `'`, `"`, -1) + `'`
}
