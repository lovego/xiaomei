package dsn

import (
	"fmt"
	"os"
	"regexp"
)

type MysqlConf struct {
	Host, Port, User, Passwd, Db string
}

func MysqlDSN(uri string) MysqlConf {
	if uri == `` {
		fmt.Println(`invalid mysql config`)
		os.Exit(1)
	}
	m := regexp.MustCompile(
		`^(\w+):(\w+)@\w+\(([^()]+):(\d+)\)/(\w+)$`,
	).FindStringSubmatch(uri)
	if len(m) == 0 {
		panic(`mysql addr match faild.`)
	}
	return MysqlConf{m[3], m[4], m[1], m[2], m[5]}
}

func (c MysqlConf) Options() []string {
	return []string{`-h` + c.Host, `-u` + c.User, `-p` + c.Passwd, `-P` + c.Port, c.Db}
}
