package cli

import (
	"fmt"
	"os"
	"regexp"

	"github.com/bughou-go/xiaomei/config"
)

type MysqlConfig struct {
	User, Passwd, Host, Port, Db string
}

func GetMysqlConfig(name string) MysqlConfig {
	dsn, ok := config.Mysql()[name]
	if !ok {
		fmt.Println(`no mysql config for: `, name)
		os.Exit(1)
	}
	m := regexp.MustCompile(
		`^(\w+):(\w+)@\w+\(([^()]+):(\d+)\)/(\w+)$`,
	).FindStringSubmatch(dsn)
	if len(m) == 0 {
		panic(`mysql addr match faild.`)
	}
	return MysqlConfig{m[1], m[2], m[3], m[4], m[5]}
}

func GetMysqlOptions(name string) []string {
	c := GetMysqlConfig(name)
	return []string{`-h` + c.Host, `-u` + c.User, `-p` + c.Passwd, `-P` + c.Port, c.Db}
}
