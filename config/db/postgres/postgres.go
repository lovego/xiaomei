package postgres

import (
	"log"
	"sync"

	"github.com/go-pg/pg"
	"github.com/lovego/xiaomei/config"
)

var postgresDBs = struct {
	sync.Mutex
	m map[string]*pg.DB
}{m: make(map[string]*pg.DB)}

func DB(name string) *pg.DB {
	postgresDBs.Lock()
	defer postgresDBs.Unlock()
	db := postgresDBs.m[name]
	if db == nil {
		db = newDB(name)
		postgresDBs.m[name] = db
	}
	return db
}

func newDB(name string) *pg.DB {
	options, err := pg.ParseURL(config.Get(`postgres`).GetString(name))
	if err != nil {
		log.Panic(err)
	}
	db := pg.Connect(options)
	return db
}
