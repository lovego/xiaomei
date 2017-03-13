package config

import (
	"time"

	"github.com/bughou-go/xiaomei/config/conf"
	"github.com/bughou-go/xiaomei/utils/mailer"
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

func Mailer() *mailer.Mailer {
	return theConf.Mailer()
}

func Alarm(title, body string) {
	theConf.Alarm(title, body)
}

func Keepers() []mailer.People {
	return theConf.Keepers()
}

func DataSource(typ, key string) string {
	return theConf.DataSource(typ, key)
}

func Parse(p interface{}) {
	conf.Parse(p, Root(), Env())
}
