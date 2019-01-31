package dbs

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/lib/pq"
	"github.com/lovego/fs"
	"github.com/lovego/xiaomei/release"
)

func setupPostgres(env string) {
	if env != `dev` && env != `test` && env != `ci` {
		return
	}
	dbUrlsMap := release.AppData(env).Get(`postgres`)
	if len(dbUrlsMap) == 0 {
		return
	}

	sqlsDir := filepath.Join(release.Root(), `../sqls`)
	for key, value := range dbUrlsMap {
		if dbUrl, ok := value.(string); ok {
			u, err := url.Parse(dbUrl)
			if err != nil {
				log.Panic(err)
			}
			dbName := strings.TrimPrefix(u.Path, `/`)
			createDB(env, dbName, u)
			file := filepath.Join(sqlsDir, key+`.sql`)
			execSqlFile(newDB(dbUrl), file)
		}
	}
}

func createDB(env, dbName string, u *url.URL) {
	u.Path = `postgres`
	_, err := newDB(u.String()).Exec(`create database ` + dbName)
	if e, ok := err.(*pq.Error); ok && e != nil && e.Get('C') == `42P04` {
		return
	}
	if err != nil {
		log.Panic(err)
	}
}

func execSqlFile(db *sql.DB, file string) {
	if fs.NotExist(file) {
		log.Printf("WARNING: file %s does not exists\n", file)
		return
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}
	sql := string(content)
	if _, err := db.Exec(sql); err != nil {
		log.Panic(err)
	}
}

func newDB(conn string) *sql.DB {
	db, err := sql.Open(`postgres`, conn)
	if err != nil {
		log.Panic(err)
	}
	return db
}
