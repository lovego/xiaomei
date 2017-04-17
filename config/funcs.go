package config

import (
	"time"

	"github.com/lovego/xiaomei/config/conf"
)

func Root() string {
	return theConf.Root()
}

func Name() string {
	return theConf.Name()
}

func Env() string {
	return theConf.Env()
}

func DeployName() string {
	return theConf.DeployName()
}

func Domain() string {
	return theConf.Domain()
}

func Secret() string {
	return theConf.Secret()
}

func TimeZone() *time.Location {
	return theConf.TimeZone()
}

func Keepers() []string {
	return theConf.Keepers()
}

func DataSource(typ, key string) string {
	return theConf.DataSource(typ, key)
}

func Parse(p interface{}) {
	conf.Parse(p, Root(), Env())
}
