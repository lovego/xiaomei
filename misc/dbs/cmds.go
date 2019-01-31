package dbs

import (
	"errors"

	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	mysqlCmd := makeCmd(`mysql`, `enter mysql cli.`, Mysql)
	mysqlCmd.AddCommand(makeCmd(`dump`, `mysqldump to stdout.`, MysqlDump), mysqlSetupCmd())

	postgresCmd := makeCmd(`psql`, `enter psql cli.`, Psql)
	postgresCmd.AddCommand(postgresSetupCmd())

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

func mysqlSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `setup <sql-file> <env> <key>`,
		Short: `create mysql databases and tables.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New(`sql-file is required.`)
			} else if len(args) > 3 {
				return errors.New(`more than three arguments given.`)
			}
			flags := c.Flags()
			file, env, key := flags.Arg(0), flags.Arg(1), flags.Arg(2)

			setupMysql(file, env, key)
			return nil
		},
	}
	return cmd
}

func postgresSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `setup <env>`,
		Short: `create postgres databases and tables.`,
		RunE: func(c *cobra.Command, args []string) error {
			env := c.Flags().Arg(0)
			if env == `` {
				env = `dev`
			}
			setupPostgres(env)
			return nil
		},
	}
	return cmd
}
