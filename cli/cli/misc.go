package cli

import (
	"regexp"
	"github.com/bughou-go/xiaomei/config"
)

type MysqlConfig struct {
	User, Passwd, Host, Port, Db string
}

var mysqlConfig MysqlConfig

func GetMysqlConfig() MysqlConfig {
	if mysqlConfig.User == `` {
		m := regexp.MustCompile(`^(\w+):(\w+)@\w+\(([^()]+):(\d+)\)/(\w+)$`).
			FindStringSubmatch(config.Data.Mysql)
		if len(m) == 0 {
			panic(`mysql addr match faild.`)
		}
		mysqlConfig = MysqlConfig{m[1], m[2], m[3], m[4], m[5]}
	}
	return mysqlConfig
}

func GetMysqlOptions() []string {
	c := GetMysqlConfig()
	return []string{`-h` + c.Host, `-u` + c.User, `-p` + c.Passwd, `-P` + c.Port, c.Db}
}
