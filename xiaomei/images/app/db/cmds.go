package db

import (
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		makeCmd(`mysql`, `enter mysql cli.`, Mysql),
		makeCmd(`mysqldump`, `mysqldump to stdout.`, MysqlDump),
		makeCmd(`mongo`, `enter mongo cli.`, Mongo),
		makeCmd(`redis`, `enter redis cli.`, Redis),
		setupCmd(),
	}
}

func makeCmd(name, short string, fun func(key string, p bool) error) *cobra.Command {
	var p bool
	cmd := &cobra.Command{
		Use:   name + ` [<dbkey>]`,
		Short: short,
		RunE: func(c *cobra.Command, args []string) error {
			return fun(c.Flags().Arg(0), p)
		},
	}
	cmd.Flags().BoolVarP(&p, `print`, `p`, false, `only print the command.`)
	return cmd
}
