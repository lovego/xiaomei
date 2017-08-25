package db

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/utils/dsn"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/release"
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
	options := dsn.Mysql(release.AppData().Get(`mysql`).GetString(key)).Flags()
	createDbAndTables(options)

	fmt.Println(`setup mysql ok.`)
}

func createDbAndTables(options []string) {
	l := len(options)
	db := options[l-1]

	createDB := fmt.Sprintf(`create database if not exists %s charset utf8; use %s;`, db, db)
	createTbs, err := ioutil.ReadFile(path.Join(release.Root(), `img-app/config/data/ddl.mysql`))
	if err != nil {
		panic(err)
	}
	sqls := bytes.NewBufferString(createDB + string(createTbs))

	cluster.Run(cmd.O{Stdin: sqls, Panic: true}, ``, `mysql `+strings.Join(options[:l-1], ` `))
}
