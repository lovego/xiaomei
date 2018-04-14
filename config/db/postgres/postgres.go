package postgres

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
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

	options.DialTimeout = 5 * time.Second
	options.ReadTimeout = 5 * time.Second
	options.WriteTimeout = 5 * time.Second
	options.IdleTimeout = time.Minute
	options.MaxAge = time.Hour
	options.PoolSize = 100

	db := pg.Connect(options)

	if os.Getenv("DebugPg") != "" {
		db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				log.Println(err)
			}
			log.Printf("Postgres: %s %s", time.Since(event.StartTime), color.GreenString(query))
		})
	}
	return db
}
