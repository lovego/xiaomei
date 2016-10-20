package setup

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/tools/tools"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func SetupMysql() {
	options := tools.GetMysqlOptions()
	createDatabaseAndTables(options)

	var v int
	err := config.Mysql().QueryRow(`select 1 from organizations limit 1`).Scan(&v)
	if err == sql.ErrNoRows {
		loadData(options) // 空表
	} else if err != nil {
		panic(err)
	}

	fmt.Println(`setup mysql ok.`)
}

func createDatabaseAndTables(options []string) {
	l := len(options)
	db := options[l-1]

	create_database := fmt.Sprintf(`create database if not exists %s charset utf8; use %s;`, db, db)
	create_tables, err := ioutil.ReadFile(path.Join(config.Root, `config/data/ddl.mysql`))
	if err != nil {
		panic(err)
	}
	sqls := bytes.NewBufferString(create_database + string(create_tables))

	cmd.Run(cmd.O{Stdin: sqls, Panic: true}, `mysql`, options[:l-1]...)
}

func loadData(options []string) {
	insert_data, err := ioutil.ReadFile(path.Join(config.Root, `config/data/data.mysql`))
	if err != nil {
		panic(err)
	}
	sqls := bytes.NewBuffer(insert_data)
	cmd.Run(cmd.O{Stdin: sqls, Panic: true}, `mysql`, options...)
}
