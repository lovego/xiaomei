package db

import (
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		makeCmd(`mysql`, `enter mysql cli`, Mysql),
		makeCmd(`mysqldump`, `dump mysql (you can redirect to a file if you need)`, MysqlDump),
		makeCmd(`mongo`, `enter mongo cli`, Mongo),
		makeCmd(`redis`, `enter redis cli`, Redis),
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
	cmd.Flags().BoolVarP(&p, `print`, `p`, false, `print the commands but do not run them.`)
	return cmd
}
