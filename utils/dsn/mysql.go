package dsn

import (
	"fmt"
	"os"
	"regexp"
)

type MysqlDSN struct {
	User, Passwd, Host, Port, Db string
}

func Mysql(uri string) MysqlDSN {
	if uri == `` {
		fmt.Println(`invalid mysql config`)
		os.Exit(1)
	}
	m := regexp.MustCompile(
		`^(\w+):(\w+)@\w+\(([^()]+):(\d+)\)/(\w+)`,
	).FindStringSubmatch(uri)
	if len(m) == 0 {
		panic(`mysql addr match faild.`)
	}
	return MysqlDSN{m[1], m[2], m[3], m[4], m[5]}
}

func (c MysqlDSN) Flags() []string {
	return []string{`-h` + c.Host, `-u` + c.User, `-p` + c.Passwd, `-P` + c.Port, c.Db}
}
