package db

import (
	"fmt"
	"strings"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/utils/dsn"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func Mysql(key string, printCmd bool) error {
	flags := dsn.Mysql(release.AppData().Get(`mysql`).GetString(key)).Flags()
	command := `mysql --pager=less -SX ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Run(cmd.O{}, ``, command)
	return err
}

func MysqlDump(key string, printCmd bool) error {
	flags := dsn.Mysql(release.AppData().Get(`mysql`).GetString(key)).Flags()
	command := `mysqldump -t ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Run(cmd.O{}, ``, command)
	return err
}

func Mongo(key string, printCmd bool) error {
	command := `mongo ` + release.AppData().Get(`mongo`).GetString(key)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Run(cmd.O{}, ``, command)
	return err
}

func Redis(key string, printCmd bool) error {
	flags := dsn.Redis(release.AppData().Get(`redis`).GetString(key)).Flags()
	command := `redis-cli ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Run(cmd.O{}, ``, command)
	return err
}
