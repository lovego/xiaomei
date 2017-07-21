package dsn

import (
	"fmt"
	"net"
	"os"

	"github.com/go-sql-driver/mysql"
)

type MysqlDSN struct {
	User, Passwd, Host, Port, Db string
}

func Mysql(uri string) MysqlDSN {
	if uri == `` {
		fmt.Println(`invalid mysql config`)
		os.Exit(1)
	}
	c, err := mysql.ParseDSN(uri)
	if err != nil {
		panic(`mysql addr match faild.`)
	}
	host, port, err := net.SplitHostPort(c.Addr)
	if err != nil {
		panic(`mysql addr match faild.`)
	}
	return MysqlDSN{c.User, c.Passwd, host, port, c.DBName}
}

func (c MysqlDSN) Flags() []string {
	return []string{`-h` + c.Host, `-u` + c.User, `-p` + c.Passwd, `-P` + c.Port, c.Db}
}
