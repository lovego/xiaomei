package db

import (
	"github.com/bughou-go/xiaomei/config"
	"net/url"
	"retail-tek.com/reports/utils/cmd"
	"strings"
)

type RedisConfig struct {
	Passwd, Host, Port, Db string
}

func Redis(key string) {
	options := getRedisOptions(key)
	command, options := sshOptions(`redis-cli`, options)
	cmd.Run(cmd.O{Panic: true}, command, options...)
}

func getRedisOptions(key string) []string {
	c := getRedisConfig(key)
	options := []string{}
	if c.Passwd != `` {
		options = append(options, `-a`, c.Passwd)
	}
	options = append(options, `-h`, c.Host, `-p`, c.Port, `-n`, c.Db)
	return options
}

func getRedisConfig(key string) RedisConfig {
	uri := config.DB.Redis(key)
	info, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	c := RedisConfig{}
	passwd, has := info.User.Password()
	if has {
		c.Passwd = passwd
	}
	hostPort := strings.Split(info.Host, `:`)
	c.Host, c.Port = hostPort[0], hostPort[1]
	c.Db = strings.Split(info.Path, `/`)[1]
	return c
}
