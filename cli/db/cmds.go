package db

import (
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		makeCmd(`mysql`, `enter mysql cli`, Mysql),
		makeCmd(`mysqldump`, `dump mysql`, MysqlDump),
		makeCmd(`mongo`, `enter mongo cli`, Mongo),
		makeCmd(`redis`, `enter redis cli`, Redis),
	}
}

func makeCmd(name, short string, fun func(key, filter string) error) *cobra.Command {
	var filter *string
	cmd := &cobra.Command{
		Use:   name + ` [<dbkey>]`,
		Short: `[db] ` + short,
		RunE: func(c *cobra.Command, args []string) error {
			return fun(c.Flags().Arg(0), *filter)
		},
	}
	filter = cmd.Flags().StringP(`server`, `s`, ``, `match servers by Addr or Tasks.`)
	return cmd
}
