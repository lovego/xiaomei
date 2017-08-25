package conf

import (
	"time"
)

type Conf struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	Secret string `yaml:"secret"`

	Mailer       string   `yaml:"mailer"`
	Keepers      []string `yaml:"keepers"`
	TimeZone     timeZone `yaml:"timeZone"`
	TimeLocation *time.Location
}

type timeZone struct {
	Name   string `yaml:"name"`
	Offset int    `yaml:"offset"`
}

func (c *Conf) DeployName(env string) string {
	return c.Name + `_` + env
}
