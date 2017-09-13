package db

import (
	"errors"

	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		makeCmd(`mysql`, `enter mysql cli.`, Mysql),
		makeCmd(`mysqldump`, `mysqldump to stdout.`, MysqlDump),
		makeCmd(`mongo`, `enter mongo cli.`, Mongo),
		makeCmd(`redis`, `enter redis cli.`, Redis),
		setupMysqlCmd(),
	}
}

func makeCmd(name, short string, fun func(env, key string, print bool) error) *cobra.Command {
	var print bool
	cmd := &cobra.Command{
		Use:   name + ` [<env> [<key>]]`,
		Short: short,
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

func setupMysqlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `setup-mysql <sql-file> <env> <key>`,
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
