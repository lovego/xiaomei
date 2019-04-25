package create

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/fatih/color"
	"github.com/lib/pq"
	"github.com/lovego/bsql"
	"github.com/lovego/config/conf"
	"github.com/lovego/xiaomei/release"
)

type Creation struct {
	env, typ, key string
	recreate      bool
	sql           string
}

func (c Creation) do() error {
	v, err := conf.GetDb(release.AppData(c.env), c.typ, c.key)
	if err != nil {
		return err
	}
	switch dbConf := v.(type) {
	case string:
		return c.doOne(dbConf, 0, conf.ShardsSettings{})
	case *conf.Shards:
		for _, shard := range dbConf.Shards {
			if err := c.doOne(shard.Url, shard.No, dbConf.Settings); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unexpected db config: %v", v)
	}
}

func (c Creation) doOne(dbUrl string, shardNo int, shardSettings conf.ShardsSettings) error {
	if dbURL, dbName, err := c.prepare(dbUrl); err != nil {
		return err
	} else if err := c.createDB(dbURL, dbName); err != nil {
		return err
	}

	db, err := sql.Open(c.typ, dbUrl)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err = db.Exec(c.sql); err != nil {
		return bsql.ErrorWithPosition(err, c.sql)
	}
	if shardNo > 0 {
		return setupShard(db, shardNo, shardSettings)
	}
	return nil
}

func (c Creation) prepare(dbUrl string) (*url.URL, string, error) {
	dbURL, err := url.Parse(dbUrl)
	if err != nil {
		return nil, "", err
	}
	dbName := strings.TrimPrefix(dbURL.Path, `/`)
	if c.recreate {
		log.Printf("recreate %s ...\n", color.GreenString(dbName))
	} else {
		log.Printf("create %s ...\n", color.GreenString(dbName))
	}
	return dbURL, dbName, nil
}

func (c Creation) createDB(dbURL *url.URL, dbName string) error {
	dbURL.Path = c.typ

	db, err := sql.Open(c.typ, dbURL.String())
	if err != nil {
		return err
	}
	defer db.Close()

	if c.recreate {
		if _, err := db.Exec(`DROP DATABASE IF EXISTS ` + dbName); err != nil {
			return err
		}
	}

	sql := `CREATE DATABASE `
	if c.typ == "mysql" {
		sql += `IF NOT EXISTS`
	}
	sql += dbName

	if _, err := db.Exec(sql); err != nil {
		if c.typ == "postgres" {
			if e, ok := err.(*pq.Error); ok && e != nil && e.Get('C') == `42P04` {
				return nil
			}
		}
		return err
	}
	return nil
}
