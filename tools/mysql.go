package main

import (
	"os"
	"path"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/tools/tools"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func mysql() {
	options := tools.GetMysqlOptions()
	options = append([]string{`--pager=less -SX`}, options...)
	cmd.Run(cmd.O{Panic: true}, `mysql`, options...)
}

func mysqldump() {
	options := tools.GetMysqlOptions()
	options = append([]string{`-t`}, options...)
	sqls, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `mysqldump`, options...)
	f, err := os.Create(path.Join(config.Root, `config/data/data.mysql`))
	if err != nil {
		panic(err)
	}
	if _, err := f.WriteString(sqls); err != nil {
		panic(err)
	}
	f.Close()
}
