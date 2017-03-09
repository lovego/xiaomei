package config

import (
	"time"

	"github.com/bughou-go/xiaomei/utils/mailer"
)

func Root() string {
	return Config.Root()
}

func Name() string {
	return Config.Name()
}

func Env() string {
	return Config.Env()
}

func DeployName() string {
	return Config.DeployName()
}

func Domain() string {
	return Config.Domain()
}

func Secret() string {
	return Config.Secret()
}

func TimeZone() *time.Location {
	return Config.TimeZone()
}

func Mailer() *mailer.Mailer {
	return Config.Mailer()
}

func Alarm(title, body string) {
	Config.Alarm(title, body)
}

func Keepers() []mailer.People {
	return Config.Keepers()
}

func DataSource(typ, key string) string {
	return Config.DataSource(typ, key)
}
