package postgres

import (
	"database/sql"
	"log"
	"net/url"
	"strconv"
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
	dataSource := newDataSource(config.Get("postgres").GetString(name))
	db, err := sql.Open("postgres", dataSource.url.String())
	if err != nil {
		log.Panic(err)
	}
	if err := db.Ping(); err != nil {
		log.Panic(err)
	}
	db.SetMaxOpenConns(dataSource.maxOpen)
	db.SetMaxIdleConns(dataSource.maxIdle)
	db.SetConnMaxLifetime(dataSource.maxLife)
	return bsql.New(db, 5*time.Second)
}

type dataSource struct {
	url              *url.URL
	maxIdle, maxOpen int
	maxLife          time.Duration
}

func newDataSource(urlStr string) dataSource {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Panic(err)
	}
	ds := dataSource{url: u}

	q := u.Query()
	if str := q.Get("maxIdle"); str != "" {
		q.Del("maxIdle")
		ds.maxIdle = parseInt(str)
	} else if config.Env() == "production" {
		ds.maxIdle = 1
	} else {
		ds.maxIdle = 0
	}

	if str := q.Get("maxOpen"); str != "" {
		q.Del("maxOpen")
		ds.maxOpen = parseInt(str)
	} else {
		ds.maxOpen = 10
	}

	if str := q.Get("maxLife"); str != "" {
		q.Del("maxLife")
		ds.maxLife = parseDuration(str)
	} else {
		ds.maxLife = 10 * time.Minute
	}

	return ds
}

func parseInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Panic(err)
	}
	return value
}

func parseDuration(str string) time.Duration {
	value, err := time.ParseDuration(str)
	if err != nil {
		log.Panic(err)
	}
	return value
}
