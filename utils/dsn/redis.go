package dsn

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type RedisDSN struct {
	Passwd, Host, Port, Db string
}

func Redis(uri string) RedisDSN {
	if uri == `` {
		fmt.Println(`invalid redis config`)
		os.Exit(1)
	}
	info, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	c := RedisDSN{}
	passwd, has := info.User.Password()
	if has {
		c.Passwd = passwd
	}
	hostPort := strings.Split(info.Host, `:`)
	c.Host, c.Port = hostPort[0], hostPort[1]
	c.Db = strings.Split(info.Path, `/`)[1]
	return c
}

func (c RedisDSN) Flags() []string {
	flags := []string{}
	if c.Passwd != `` {
		flags = append(flags, `-a`, c.Passwd)
	}
	flags = append(flags, `-h`, c.Host, `-p`, c.Port, `-n`, c.Db)
	return flags
}
