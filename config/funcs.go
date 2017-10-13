package config

import (
	"time"

	"github.com/lovego/xiaomei/utils/mailer"
	"github.com/lovego/xiaomei/utils/strmap"
)

func Name() string {
	return theConf.Name
}

func DeployName() string {
	return theConf.DeployName()
}

func Domain() string {
	return theConf.Domain
}

func Secret() string {
	return theConf.Secret
}

func TimeZone() *time.Location {
	return theConf.TimeLocation
}

func Keepers() []string {
	return theConf.Keepers
}

func Mailer() *mailer.Mailer {
	return theMailer
}

func Get(key string) strmap.StrMap {
	return theData.Get(key)
}

func GetSlice(key string) []strmap.StrMap {
	return theData.GetSlice(key)
}

func GetString(key string) string {
	return theData.GetString(key)
}

func GetStringSlice(key string) []string {
	return theData.GetStringSlice(key)
}
