package dbs

import (
	"fmt"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/config/db/dburl"
	"github.com/lovego/dsn"
	"github.com/lovego/xiaomei/cluster"
	"github.com/lovego/xiaomei/release"
)

func Psql(env, key string, printCmd bool) error {
	url := dburl.Parse(release.AppData(env).Get(`postgres`).GetString(key)).URL.String()
	command := fmt.Sprintf("PAGER=less LESS='-iFMSXx4' psql '%s'", url)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Get(env).Run(``, cmd.O{}, command)
	return err
}

func Mysql(env, key string, printCmd bool) error {
	flags := dsn.Mysql(release.AppData(env).Get(`mysql`).GetString(key)).Flags()
	command := `mysql --pager=less -SX ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Get(env).Run(``, cmd.O{}, command)
	return err
}

func MysqlDump(env, key string, printCmd bool) error {
	flags := dsn.Mysql(release.AppData(env).Get(`mysql`).GetString(key)).Flags()
	command := `mysqldump -t ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Get(env).Run(``, cmd.O{}, command)
	return err
}

func Mongo(env, key string, printCmd bool) error {
	command := `mongo ` + release.AppData(env).Get(`mongo`).GetString(key)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Get(env).Run(``, cmd.O{}, command)
	return err
}

func Redis(env, key string, printCmd bool) error {
	flags := dsn.Redis(release.AppData(env).Get(`redis`).GetString(key)).Flags()
	command := `redis-cli ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := cluster.Get(env).Run(``, cmd.O{}, command)
	return err
}
