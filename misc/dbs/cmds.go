package dbs

import (
	"strings"

	"github.com/lovego/xiaomei/misc/dbs/create"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	mysqlCmd := makeCmd(`mysql`, `Enter mysql cli.`, Mysql)
	mysqlCmd.AddCommand(
		makeCmd(`dump`, `mysqldump to stdout.`, MysqlDump),
	)

	postgresCmd := makeCmd(`psql`, `Enter psql cli.`, Psql)
	postgresCmd.AddCommand(
		makeCreateCmd(`postgres`, `create`, ``),
		makeCreateCmd(`postgres`, `recreate`, `drop existing databases and `),
	)

	return []*cobra.Command{
		postgresCmd,
		mysqlCmd,
		makeCmd(`mongo`, `Enter mongo cli.`, Mongo),
		makeCmd(`redis`, `Enter redis cli.`, Redis),
	}
}

func makeCmd(name, short string, fun func(env, key string, print bool) error) *cobra.Command {
	var print bool
	cmd := &cobra.Command{
		Use:   name + ` [<env> [<key>]]`,
		Short: `[db] ` + short,
		RunE: release.Env1Call(func(env, key string) error {
			if key == `` {
				key = `default`
			}
			return fun(env, key, print)
		}),
	}
	cmd.Flags().BoolVarP(&print, `print`, `p`, false, `only print the command.`)
	return cmd
}

func makeCreateCmd(typ, name, short string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name + ` [<env> [<key>]]`,
		Short: strings.Title(name) + ` databases and tables (` + short + `execute ".sql" file in "sqls" dir).`,
		RunE: release.Env1Call(func(env, key string) error {
			return create.Do(env, typ, key, name == `recreate`)
		}),
	}
	return cmd
}
