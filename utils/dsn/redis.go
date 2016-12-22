package dsn

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type RedisConf struct {
	Passwd, Host, Port, Db string
}

func RedisDSN(uri string) RedisConf {
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

func (c RedisConf) Options() []string {
	options := []string{}
	if c.Passwd != `` {
		options = append(options, `-a`, c.Passwd)
	}
	options = append(options, `-h`, c.Host, `-p`, c.Port, `-n`, c.Db)
	return options
}
