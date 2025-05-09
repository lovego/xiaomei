package dbs

import (
	"fmt"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/config/db/dburl"
	"github.com/lovego/xiaomei/release"
)

const less = "'less -iFMSXx4'"

func Psql(env, key string, lessPager, printCmd bool) error {
	url := dburl.Parse(getDBUrl(env, `postgres`, key)).URL.String()
	command := fmt.Sprintf("psql '%s'", url)
	if lessPager {
		command = fmt.Sprintf("PAGER=%s %s", less, command)
	}
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := release.GetCluster(env).Run(``, cmd.O{}, command)
	return err
}

func Mysql(env, key string, lessPager, printCmd bool) error {
	command := `mysql `
	if lessPager {
		command += `--pager=` + less + " "
	}
	command += strings.Join(ParseDSN(getDBUrl(env, `mysql`, key)).MysqlFlags(), ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := release.GetCluster(env).Run(``, cmd.O{}, command)
	return err
}

func MysqlDump(env, key string, printCmd bool) error {
	flags := ParseDSN(getDBUrl(env, `mysql`, key)).MysqlFlags()
	command := `mysqldump -t ` + strings.Join(flags, ` `)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := release.GetCluster(env).Run(``, cmd.O{}, command)
	return err
}

func Mongo(env, key string, printCmd bool) error {
	command := `mongo ` + getDBUrl(env, `mongo`, key)
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := release.GetCluster(env).Run(``, cmd.O{}, command)
	return err
}

func Redis(env, key string, printCmd bool) error {
	command := `redis-cli -u ` + strings.Split(getDBUrl(env, `redis`, key), "?")[0]
	if printCmd {
		fmt.Println(command)
		return nil
	}
	_, err := release.GetCluster(env).Run(``, cmd.O{}, command)
	return err
}

func getDBUrl(env, typ, key string) string {
	strMap := release.Config(env).Data.Get(typ)
	keys := strings.Split(key, ".")
	for i := 0; i < len(keys)-1; i++ {
		strMap = strMap.Get(keys[i])
	}
	return strMap.GetString(keys[len(keys)-1])
}
