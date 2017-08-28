package db

import (
	"errors"
	"fmt"

	"github.com/lovego/xiaomei/utils/slice"
	"github.com/lovego/xiaomei/xiaomei/cluster"
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
		Use:   name + ` [<env> [<dbkey>]]`,
		Short: short,
		RunE: func(c *cobra.Command, args []string) error {
			env, key, err := GetEnvAndKey(args)
			if err != nil {
				return err
			}
			return fun(env, key, print)
		},
	}
	cmd.Flags().BoolVarP(&print, `print`, `p`, false, `only print the command.`)
	return cmd
}

func setupMysqlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `setup-mysql [<env> [<dbkey>]]`,
		Short: `create mysql databases and tables.`,
		RunE: func(c *cobra.Command, args []string) error {
			env, key, err := GetEnvAndKey(args)
			if err != nil {
				return err
			}
			setupMysql(env, key)
			return nil
		},
	}
	return cmd
}

func GetEnvAndKey(args []string) (env, key string, err error) {
	switch len(args) {
	case 0:
	case 1:
		env = args[0]
	case 2:
		env = args[0]
		key = args[1]
	default:
		err = errors.New(`more than two arguments given.`)
		return
	}
	if env == `` {
		env = `dev`
	}
	if key == `` {
		key = `default`
	}
	if !slice.ContainsString(cluster.Envs(), env) {
		err = fmt.Errorf("env %s not defined in cluster.yml", env)
	}
	return
}
