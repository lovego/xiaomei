package utils

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type MysqlConf struct {
	Host, Port, User, Passwd, Db string
}

func MysqlConfig(dsn string) MysqlConf {
	if dsn == `` {
		fmt.Println(`invalid mysql config`)
		os.Exit(1)
	}
	m := regexp.MustCompile(
		`^(\w+):(\w+)@\w+\(([^()]+):(\d+)\)/(\w+)$`,
	).FindStringSubmatch(dsn)
	if len(m) == 0 {
		panic(`mysql addr match faild.`)
	}
	return MysqlConf{m[3], m[4], m[1], m[2], m[5]}
}

func MysqlOptions(uri string) []string {
	c := MysqlConfig(uri)
	return []string{`-h` + c.Host, `-u` + c.User, `-p` + c.Passwd, `-P` + c.Port, c.Db}
}

type RedisConf struct {
	Passwd, Host, Port, Db string
}

func RedisOptions(uri string) []string {
	c := RedisConfig(uri)
	options := []string{}
	if c.Passwd != `` {
		options = append(options, `-a`, c.Passwd)
	}
	options = append(options, `-h`, c.Host, `-p`, c.Port, `-n`, c.Db)
	return options
}

func RedisConfig(uri string) RedisConf {
	if uri == `` {
		fmt.Println(`invalid redis config`)
		os.Exit(1)
	}
	info, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	c := RedisConf{}
	passwd, has := info.User.Password()
	if has {
		c.Passwd = passwd
	}
	hostPort := strings.Split(info.Host, `:`)
	c.Host, c.Port = hostPort[0], hostPort[1]
	c.Db = strings.Split(info.Path, `/`)[1]
	return c
}
