package dbs

import (
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	mysqlCmd := makeCmd(`mysql`, `enter mysql cli.`, Mysql)
	mysqlCmd.AddCommand(
		makeCmd(`dump`, `mysqldump to stdout.`, MysqlDump),
	)

	postgresCmd := makeCmd(`psql`, `enter psql cli.`, Psql)
	postgresCmd.AddCommand(makeSetupCmd(`postgres`))

	return []*cobra.Command{
		postgresCmd,
		mysqlCmd,
		makeCmd(`mongo`, `enter mongo cli.`, Mongo),
		makeCmd(`redis`, `enter redis cli.`, Redis),
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

func makeSetupCmd(typ string) *cobra.Command {
	var drop bool
	cmd := &cobra.Command{
		Use:   `setup [<env> [<key>]]`,
		Short: `setup databases and tables. (execute .sql file in "sqls" dir).`,
		RunE: release.Env1Call(func(env, key string) error {
			return setup(env, typ, key, drop)
		}),
	}
	cmd.Flags().BoolVarP(&drop, `drop`, `d`, false, `drop existing database before setup.`)
	return cmd
}
