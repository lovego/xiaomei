package setup

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/dsn"
)

func SetupMysql() {
	options := dsn.MysqlDSN(config.DB.Mysql(``)).Options()
	if msg := createDatabaseAndTables(options); msg == `ok` {
		fmt.Println(`setup mysql ok.`)
	} else {
		fmt.Println(msg)
	}
}

func createDatabaseAndTables(options []string) string {
	l := len(options)
	db := options[l-1]

	createDb := fmt.Sprintf(`create database if not exists %s charset utf8; use %s;`, db, db)
	filePath := path.Join(config.App.Root(), `config/data/ddl.mysql`)
	if !utils.IsFile(filePath) {
		return `no such file: ` + filePath
	}
	createTables, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	sqls := bytes.NewBufferString(createDb + string(createTables))

	cmd.Run(cmd.O{Stdin: sqls, Panic: true}, `mysql`, options[:l-1]...)
	return `ok`
}

func loadData(options []string) {
	insert_data, err := ioutil.ReadFile(path.Join(config.App.Root(), `config/data/data.mysql`))
	if err != nil {
		panic(err)
	}
	sqls := bytes.NewBuffer(insert_data)
	cmd.Run(cmd.O{Stdin: sqls, Panic: true}, `mysql`, options...)
}
