package dbs

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/lib/pq"
	"github.com/lovego/bsql"
	"github.com/lovego/strmap"
	"github.com/lovego/xiaomei/release"
)

func create(env, typ, key string, recreate bool) error {
	if recreate {
		switch env {
		case `production`, `staging`, `preview`:
			return errors.New("recreate is forbidden under environments: " + env)
		}
	}

	var keys []string
	if key != "" {
		keys = []string{key}
	} else {
		files, err := filepath.Glob(filepath.Join(release.Root(), `..`, `sqls`, `*.sql`))
		if err != nil {
			return err
		}
		for _, file := range files {
			keys = append(keys, strings.TrimSuffix(filepath.Base(file), `.sql`))
		}
	}

	dbConfig := release.AppData(env).Get(typ)
	for _, key := range keys {
		if err := createDatabases(typ, key, dbConfig, recreate); err != nil {
			return err
		}
	}
	return nil
}

func createDatabases(typ, key string, dbConfig strmap.StrMap, recreate bool) error {
	dbUrls, err := getDbUrls(dbConfig, key)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(filepath.Join(release.Root(), `../sqls`, key+`.sql`))
	if err != nil {
		return err
	}
	sqlContent := string(content)
	for _, dbUrl := range dbUrls {
		if err := createDatabase(typ, dbUrl, recreate, sqlContent); err != nil {
			return err
		}
	}
	return nil
}

func createDatabase(typ, dbUrl string, recreate bool, sqlContent string) error {
	dbURL, err := url.Parse(dbUrl)
	if err != nil {
		return err
	}
	dbName := strings.TrimPrefix(dbURL.Path, `/`)
	if recreate {
		log.Printf("recreate %s ...\n", color.GreenString(dbName))
	} else {
		log.Printf("create %s ...\n", color.GreenString(dbName))
	}

	if err := createDB(typ, dbURL, dbName, recreate); err != nil {
		return err
	}
	return execSqlFile(typ, dbUrl, sqlContent)
}

func createDB(typ string, dbURL *url.URL, dbName string, recreate bool) error {
	dbURL.Path = typ

	db, err := sql.Open(typ, dbURL.String())
	if err != nil {
		return err
	}
	defer db.Close()

	if recreate {
		if _, err := db.Exec(`DROP DATABASE IF EXISTS ` + dbName); err != nil {
			return err
		}
	}

	sql := `CREATE DATABASE `
	if typ == "mysql" {
		sql += `IF NOT EXISTS`
	}
	sql += dbName

	if _, err := db.Exec(sql); err != nil {
		if typ == "postgres" {
			if e, ok := err.(*pq.Error); ok && e != nil && e.Get('C') == `42P04` {
				return nil
			}
		}
		return err
	}
	return nil
}

func execSqlFile(typ, dbUrl, sqlContent string) error {
	db, err := sql.Open(typ, dbUrl)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(sqlContent)
	return bsql.ErrorWithPosition(err, sqlContent)
}

func getDbUrls(dbConfig strmap.StrMap, key string) ([]string, error) {
	value, ok := dbConfig[key]
	if !ok {
		return nil, errors.New("no database url config for: " + key)
	}
	switch v := value.(type) {
	case string:
		return []string{v}, nil
	case map[interface{}]interface{}:
		return getShardsDbUrls(v)
	default:
		return nil, errors.New("no database url config for: " + key)
	}
}

func getShardsDbUrls(m map[interface{}]interface{}) ([]string, error) {
	shardKeys := make([]int, 0, len(m))
	for k, _ := range m {
		if shard, ok := k.(int); ok {
			shardKeys = append(shardKeys, shard)
		} else {
			return nil, fmt.Errorf("shard key should be a integer, but got: %v", k)
		}
	}
	sort.Ints(shardKeys)

	dbUrls := make([]string, 0, len(m))
	for _, shard := range shardKeys {
		if dbUrl, ok := m[shard].(string); ok {
			dbUrls = append(dbUrls, dbUrl)
		} else {
			return nil, fmt.Errorf("shard db url should be a string, but got: %v", m[shard])
		}
	}
	return dbUrls, nil
}
