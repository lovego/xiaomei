package conf

import (
	"sync"
	"time"
)

type Conf struct {
	sync.Mutex

	root, env string
	data      *conf
	timeZone  *time.Location
}

type conf struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	Secret string `yaml:"secret"`

	TimeZone TimeZoneConf `yaml:"timeZone"`
	Keepers  []string     `yaml:"keepers"`

	DataSource map[string]map[string]string `yaml:"dataSource"`
}

type TimeZoneConf struct {
	Name   string `yaml:"name"`
	Offset int    `yaml:"offset"`
}

func (c *Conf) Root() string {
	return c.root
}

func (c *Conf) Env() string {
	return c.env
}

func (c *Conf) Name() string {
	return c.data.Name
}

func (c *Conf) DeployName() string {
	return c.Name() + `_` + c.Env()
}

func (c *Conf) Domain() string {
	return c.data.Domain
}

func (c *Conf) Secret() string {
	return c.data.Secret
}

func (c *Conf) TimeZone() *time.Location {
	if c.timeZone == nil {
		c.timeZone = time.FixedZone(c.data.TimeZone.Name, c.data.TimeZone.Offset)
	}
	return c.timeZone
}

func (c *Conf) Keepers() []string {
	return c.data.Keepers
}

func (c *Conf) DataSource(typ, key string) string {
	if key == `` {
		key = `default`
	}
	return c.data.DataSource[typ][key]
}
