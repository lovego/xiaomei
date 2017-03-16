package db

import (
	"fmt"
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/dsn"
	"github.com/bughou-go/xiaomei/xiaomei/cluster"
	"github.com/bughou-go/xiaomei/xiaomei/release"
)

func Mysql(key string, printCmd bool) error {
	flags := dsn.Mysql(release.App().DataSource(`mysql`, key)).Flags()
	command := `mysql --pager=less -SX ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Run(cmd.O{}, command)
	return err
}

func MysqlDump(key string, printCmd bool) error {
	flags := dsn.Mysql(release.App().DataSource(`mysql`, key)).Flags()
	command := `mysqldump -t ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Run(cmd.O{}, command)
	return err
}

func Mongo(key string, printCmd bool) error {
	command := `mongo ` + release.App().DataSource(`mongo`, key)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Run(cmd.O{}, command)
	return err
}

func Redis(key string, printCmd bool) error {
	flags := dsn.Redis(release.App().DataSource(`redis`, key)).Flags()
	command := `redis-cli ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Run(cmd.O{}, command)
	return err
}
