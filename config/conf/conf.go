package conf

import (
	"log"
	"sync"
	"time"

	"github.com/lovego/xiaomei/utils/mailer"
)

type Conf struct {
	sync.Mutex

	root, env string
	data      *conf
	timeZone  *time.Location
	mailer    struct {
		setted bool
		*mailer.Mailer
	}
}

type conf struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	Secret string `yaml:"secret"`

	TimeZone TimeZoneConf `yaml:"timeZone"`
	Mailer   MailerConf   `yaml:"mailer"`
	Keepers  []string     `yaml:"keepers"`

	DataSource map[string]map[string]string `yaml:"dataSource"`
}

type TimeZoneConf struct {
	Name   string `yaml:"name"`
	Offset int    `yaml:"offset"`
}

type MailerConf struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Sender string `yaml:"sender"`
	Passwd string `yaml:"passwd"`
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

func (c *Conf) Mailer() *mailer.Mailer {
	c.Lock()
	defer c.Unlock()
	if !c.mailer.setted {
		m := c.data.Mailer
		mail, err := mailer.New(m.Host, m.Port, m.Passwd, m.Sender)
		if err != nil {
			log.Println(err)
		}
		c.mailer.Mailer = mail
		c.mailer.setted = true
	}
	return c.mailer.Mailer
}

func (c *Conf) Alarm(title, body string) {
	title = c.DeployName() + ` ` + title
	msg := c.Mailer().NewMessage(c.Keepers(), nil, title, body, ``)
	c.Mailer().Send(msg)
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
