package setup

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/dsn"
	"github.com/bughou-go/xiaomei/utils/fs"
)

func SetupMysql() {
	flags := dsn.Mysql(config.DB.Mysql(``)).Flags()
	if msg := createDatabaseAndTables(flags); msg == `ok` {
		fmt.Println(`setup mysql ok.`)
	} else {
		fmt.Println(msg)
	}
}

func createDatabaseAndTables(flags []string) string {
	l := len(flags)
	db := flags[l-1]

	createDb := fmt.Sprintf(`create database if not exists %s charset utf8; use %s;`, db, db)
	filePath := path.Join(config.App.Root(), `config/data/ddl.mysql`)
	if !fs.IsFile(filePath) {
		return `no such file: ` + filePath
	}
	createTables, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	sqls := bytes.NewBufferString(createDb + string(createTables))

	cmd.Run(cmd.O{Stdin: sqls, Panic: true}, `mysql`, flags[:l-1]...)
	return `ok`
}

func loadData(flags []string) {
	insert_data, err := ioutil.ReadFile(path.Join(config.App.Root(), `config/data/data.mysql`))
	if err != nil {
		panic(err)
	}
	sqls := bytes.NewBuffer(insert_data)
	cmd.Run(cmd.O{Stdin: sqls, Panic: true}, `mysql`, flags...)
}
