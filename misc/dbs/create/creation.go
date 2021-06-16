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
	"github.com/lovego/config/config"
	"github.com/lovego/config/db/dburl"
	"github.com/lovego/xiaomei/release"
)

type Creation struct {
	env, typ, key string
	dropDB        bool
	sql           string
}

func (c Creation) do() error {
	v, err := config.GetDB(release.EnvData(c.env), c.typ, c.key)
	if err != nil {
		return err
	}
	switch dbConf := v.(type) {
	case string:
		return c.doOne(dbConf, 0, config.ShardsSettings{})
	case *config.Shards:
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

func (c Creation) doOne(dbUrl string, shardNo int, shardSettings config.ShardsSettings) error {
	dbURL := dburl.Parse(dbUrl).URL
	dbName := strings.TrimPrefix(dbURL.Path, `/`)
	// make a dbURL copy
	if err := c.createDB(*dbURL, dbName); err != nil {
		return err
	}

	db, err := sql.Open(c.typ, dbURL.String())
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err = db.Exec(c.sql); err != nil {
		return bsql.ErrorWithPosition(err, c.sql, false)
	}
	if shardNo > 0 {
		return c.setupShard(db, shardNo, shardSettings)
	}
	return nil
}

func (c Creation) createDB(dbURL url.URL, dbName string) error {
	dbURL.Path = c.typ

	db, err := sql.Open(c.typ, dbURL.String())
	if err != nil {
		return err
	}
	defer db.Close()

	if c.dropDB {
		log.Printf("drop database %s\n", color.GreenString(dbName))
		if _, err := db.Exec(`DROP DATABASE IF EXISTS ` + dbName); err != nil {
			return err
		}
	}

	log.Printf("create %s ...\n", color.GreenString(dbName))

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
