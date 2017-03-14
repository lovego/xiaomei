package db

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/dsn"
	"github.com/bughou-go/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func setupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `setup-mysql [<dbkey>]`,
		Short: `create mysql databases and tables.`,
		Run: func(c *cobra.Command, args []string) {
			key := ``
			if len(args) > 0 {
				key = args[0]
			}
			setupMysql(key)
		},
	}
	return cmd
}

func setupMysql(key string) {
	options := dsn.Mysql(release.App().DataSource(`mysql`, key)).Flags()
	createDbAndTables(options)

	fmt.Println(`setup mysql ok.`)
}

func createDbAndTables(options []string) {
	l := len(options)
	db := options[l-1]

	createDB := fmt.Sprintf(`create database if not exists %s charset utf8; use %s;`, db, db)
	createTbs, err := ioutil.ReadFile(path.Join(release.App().Root(), `config/data/ddl.mysql`))
	if err != nil {
		panic(err)
	}
	sqls := bytes.NewBufferString(createDB + string(createTbs))

	cmd.Run(cmd.O{Stdin: sqls, Panic: true}, `mysql`, options[:l-1]...)
}
