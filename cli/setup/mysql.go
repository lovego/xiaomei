package setup

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func SetupMysql() {
	options := config.DB.MysqlOptions(``)
	createDatabaseAndTables(options)

	fmt.Println(`setup mysql ok.`)
}

func createDatabaseAndTables(options []string) {
	l := len(options)
	db := options[l-1]

	createDb := fmt.Sprintf(`create database if not exists %s charset utf8; use %s;`, db, db)
	createTables, err := ioutil.ReadFile(path.Join(config.App.Root(), `config/data/ddl.mysql`))
	if err != nil {
		panic(err)
	}
	sqls := bytes.NewBufferString(createDb + string(createTables))

	cmd.Run(cmd.O{Stdin: sqls, Panic: true}, `mysql`, options[:l-1]...)
}

func loadData(options []string) {
	insert_data, err := ioutil.ReadFile(path.Join(config.App.Root(), `config/data/data.mysql`))
	if err != nil {
		panic(err)
	}
	sqls := bytes.NewBuffer(insert_data)
	cmd.Run(cmd.O{Stdin: sqls, Panic: true}, `mysql`, options...)
}
