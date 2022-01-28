package dbs

import (
	"github.com/lovego/xiaomei/misc/dbs/create"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	mysqlCmd := makeCmdWithLess(`mysql`, `Enter mysql cli.`, Mysql)
	mysqlCmd.AddCommand(
		makeCmd(`dump`, `mysqldump to stdout.`, MysqlDump),
	)

	postgresCmd := makeCmdWithLess(`psql`, `Enter psql cli.`, Psql)
	postgresCmd.AddCommand(makeCreateCmd(`postgres`))

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
		Use:   name + ` [env [key]]`,
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

func makeCmdWithLess(name, short string, fun func(env, key string, less, print bool) error) *cobra.Command {
	var less bool
	var print bool
	cmd := &cobra.Command{
		Use:   name + ` [env [key]]`,
		Short: `[db] ` + short,
		RunE: release.Env1Call(func(env, key string) error {
			if key == `` {
				key = `default`
			}
			return fun(env, key, less, print)
		}),
	}
	cmd.Flags().BoolVarP(&less, `less`, `l`, false, `use less as the pager.`)
	cmd.Flags().BoolVarP(&print, `print`, `p`, false, `only print the command.`)
	return cmd
}

func makeCreateCmd(typ string) *cobra.Command {
	var dropDB bool
	cmd := &cobra.Command{
		Use:   `create [env [key]]`,
		Short: `create databases and tables (execute "*.sql" file in "sqls" dir).`,
		RunE: release.Env1Call(func(env, key string) error {
			return create.Do(env, typ, key, dropDB)
		}),
	}
	cmd.Flags().BoolVarP(&dropDB, "drop-db", "d", false, "drop existing databases before creating")
	return cmd
}
