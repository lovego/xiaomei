package dbs

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lib/pq"
	"github.com/lovego/xiaomei/release"
)

func setup(env, typ, key string, dropDatabase bool) error {
	if env != `dev` && env != `test` && env != `ci` {
		return errors.New("not allowed env: " + env)
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

	dbUrls := release.AppData(env).Get(typ)
	for _, key := range keys {
		if err := setupDatabase(typ, key, dbUrls.GetString(key), dropDatabase); err != nil {
			return err
		}
	}
	return nil
}

func setupDatabase(typ, key, dbUrl string, dropDatabase bool) error {
	dbURL, err := url.Parse(dbUrl)
	if err != nil {
		return err
	}
	dbName := strings.TrimPrefix(dbURL.Path, `/`)
	log.Printf("setup %s ...\n", color.GreenString(dbName))

	if err := createDB(typ, dbURL, dbName, dropDatabase); err != nil {
		return err
	}
	return execSqlFile(typ, dbUrl, key)
}

func createDB(typ string, dbURL *url.URL, dbName string, dropDatabase bool) error {
	dbURL.Path = typ

	db, err := sql.Open(typ, dbURL.String())
	if err != nil {
		return err
	}
	defer db.Close()

	if dropDatabase {
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

func execSqlFile(typ, dbUrl string, key string) error {
	content, err := ioutil.ReadFile(filepath.Join(release.Root(), `../sqls`, key+`.sql`))
	if err != nil {
		return err
	}

	db, err := sql.Open(typ, dbUrl)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(string(content))
	return err
}
