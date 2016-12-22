package db

import (
	"os"
	"path"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/dsn"
)

func Mysql(name string) {
	options := dsn.MysqlDSN(config.DB.Mysql(name)).Options()
	options = append([]string{`--pager=less -SX`}, options...)
	cmd.Run(cmd.O{Panic: true}, `mysql`, options...)
}

func Mysqldump(name string) {
	options := dsn.MysqlDSN(config.DB.Mysql(name)).Options()
	options = append([]string{`-t`}, options...)
	sqls, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `mysqldump`, options...)
	f, err := os.Create(path.Join(config.App.Root(), `config/data/data.mysql`))
	if err != nil {
		panic(err)
	}
	if _, err := f.WriteString(sqls); err != nil {
		panic(err)
	}
	f.Close()
}
