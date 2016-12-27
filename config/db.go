package config

import ()

var DB DbConf

type DbConf struct {
	conf dbConf
}

type dbConf struct {
	Mysql map[string]string `yaml:"mysql"`
	Redis map[string]string `yaml:"redis"`
	Mongo map[string]string `yaml:"mongo"`
}

func (db *DbConf) Redis(name string) string {
	Load()
	if name == `` {
		name = `default`
	}
	return db.conf.Redis[name]
}

func (db *DbConf) Mysql(name string) string {
	Load()
	if name == `` {
		name = `default`
	}
	return db.conf.Mysql[name]
}

func (db *DbConf) Mongo(name string) string {
	Load()
	if name == `` {
		name = "default"
	}
	return db.conf.Mongo[name]
}
