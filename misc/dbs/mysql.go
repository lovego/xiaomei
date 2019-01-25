package dbs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/dsn"
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/release/cluster"
)

func setupMysql(file, env, key string) {
	options := dsn.Mysql(release.AppData(env).Get(`mysql`).GetString(key)).Flags()
	createDbAndTables(file, env, options)

	fmt.Println(`setup mysql ok.`)
}

func createDbAndTables(file, env string, options []string) {
	l := len(options)
	db := options[l-1]

	createDB := fmt.Sprintf(`create database if not exists %s charset utf8; use %s;`, db, db)
	createTbs, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	sqls := bytes.NewBufferString(createDB + string(createTbs))

	cluster.Get(env).Run(``, cmd.O{Stdin: sqls, Panic: true}, `mysql `+strings.Join(options[:l-1], ` `))
}
