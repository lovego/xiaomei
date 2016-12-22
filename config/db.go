package config

import ()

var DB dbVar

type dbVar struct {
	conf dbConf
}

type dbConf struct {
	Mysql map[string]string `yaml:"mysql"`
	Redis map[string]string `yaml:"redis"`
	Mongo map[string]string `yaml:"mongo"`
}

func (db *dbVar) Redis(name string) string {
	Load()
	if name == `` {
		name = `default`
	}
	return db.conf.Redis[name]
}

func (db *dbVar) Mysql(name string) string {
	Load()
	if name == `` {
		name = `default`
	}
	return db.conf.Mysql[name]
}

func (db *dbVar) Mongo(name string) string {
	Load()
	if name == `` {
		name = "default"
	}
	return db.conf.Mongo[name]
}
