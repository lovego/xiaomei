package dbs

import (
	"fmt"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/config/db/dburl"
	"github.com/lovego/dsn"
	"github.com/lovego/xiaomei/release"
)

func Psql(env, key string, printCmd bool) error {
	url := dburl.Parse(getDbUrl(env, `postgres`, key)).URL.String()
	command := fmt.Sprintf("PAGER='less -iFMSXx4' psql '%s'", url)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := release.GetCluster(env).Run(``, cmd.O{}, command)
	return err
}

func Mysql(env, key string, printCmd bool) error {
	flags := dsn.Mysql(getDbUrl(env, `mysql`, key)).Flags()
	command := `mysql --pager=less -SX ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := release.GetCluster(env).Run(``, cmd.O{}, command)
	return err
}

func MysqlDump(env, key string, printCmd bool) error {
	flags := dsn.Mysql(getDbUrl(env, `mysql`, key)).Flags()
	command := `mysqldump -t ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := release.GetCluster(env).Run(``, cmd.O{}, command)
	return err
}

func Mongo(env, key string, printCmd bool) error {
	command := `mongo ` + getDbUrl(env, `mongo`, key)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := release.GetCluster(env).Run(``, cmd.O{}, command)
	return err
}

func Redis(env, key string, printCmd bool) error {
	flags := dsn.Redis(getDbUrl(env, `redis`, key)).Flags()
	command := `redis-cli ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := release.GetCluster(env).Run(``, cmd.O{}, command)
	return err
}

func getDbUrl(env, typ, key string) string {
	strMap := release.AppData(env).Get(typ)
	keys := strings.Split(key, ".")
	for i := 0; i < len(keys)-1; i++ {
		strMap = strMap.Get(keys[i])
	}
	return strMap.GetString(keys[len(keys)-1])
}
