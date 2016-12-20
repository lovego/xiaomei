package config

import (
	"fmt"
	"os"
	"regexp"
)

var DB dbVar

type dbVar struct {
	conf dbConf
}

type dbConf struct {
	Mysql map[string]string `yaml:"mysql"`
	Redis map[string]string `yaml:"redis"`
}

func (db *dbVar) Redis(name string) string {
	if name == `` {
		name = `default`
	}
	return db.conf.Redis[name]
}

func (db *dbVar) Mysql(name string) string {
	if name == `` {
		name = `default`
	}
	return db.conf.Mysql[name]
}

type MysqlConf struct {
	Host, Port, User, Passwd, Db string
}

func (db *dbVar) MysqlConfig(name string) MysqlConf {
	dsn := db.Mysql(name)
	if dsn == `` {
		fmt.Println(`no mysql config for: `, name)
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

func (db *dbVar) MysqlOptions(name string) []string {
	c := db.MysqlConfig(name)
	return []string{`-h` + c.Host, `-u` + c.User, `-p` + c.Passwd, `-P` + c.Port, c.Db}
}
