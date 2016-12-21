package config

import (
	"sync"
	"time"

	"github.com/bughou-go/xiaomei/utils/mailer"
)

var App appVar

type appVar struct {
	root         string
	conf         appConf
	startTimeout struct {
		setted bool
		time.Duration
	}
	timeZone *time.Location
	mailer   struct {
		sync.Mutex
		setted bool
		*mailer.Mailer
	}
}

type appConf struct {
	Name         string `yaml:"name"`
	Env          string `yaml:"env"`
	Port         string `yaml:"port"`
	Domain       string `yaml:"domain"`
	Secret       string `yaml:"secret"`
	StartTimeout string `yaml:"startTimeout"`

	TimeZone TimeZoneConf `yaml:"timeZone"`
	Mailer   MailerConf   `yaml:"mailer"`
	Keeper   []string     `yaml:"keeper"`
}

type TimeZoneConf struct {
	Name   string `yaml:"name"`
	Offset int    `yaml:"offset"`
}

type MailerConf struct {
	Host, Port, Sender, Passwd string
}

func (a *appVar) Root() string {
	if a.root == `` {
		if root := detectRoot(); root != `` {
			a.root = root
		} else {
			panic(`app root not found.`)
		}
	}
	return a.root
}

func (a *appVar) Name() string {
	Load()
	return a.conf.Name
}

func (a *appVar) Port() string {
	Load()
	return a.conf.Port
}

func (a *appVar) Env() string {
	Load()
	return a.conf.Env
}

func (a *appVar) Domain() string {
	Load()
	return a.conf.Domain
}

func (a *appVar) Secret() string {
	Load()
	return a.conf.Secret
}

func (a *appVar) StartTimeout() time.Duration {
	if !a.startTimeout.setted {
		Load()
		if d, err := time.ParseDuration(a.conf.StartTimeout); err != nil {
			panic(err)
		} else {
			a.startTimeout.Duration = d
			a.startTimeout.setted = true
		}
	}
	return a.startTimeout.Duration
}

func (a *appVar) TimeZone() *time.Location {
	if a.timeZone == nil {
		Load()
		a.timeZone = time.FixedZone(a.conf.TimeZone.Name, a.conf.TimeZone.Offset)
	}
	return a.timeZone
}

func (a *appVar) Mailer() *mailer.Mailer {
	a.mailer.Lock()
	defer a.mailer.Unlock()
	if !a.mailer.setted {
		Load()
		m := a.conf.Mailer
		if m.Host != `` && m.Port != `` && m.Sender != `` {
			a.mailer.Mailer = mailer.New(m.Host, m.Port, m.Sender, m.Passwd)
		}
		a.mailer.setted = true
	}
	return a.mailer.Mailer
}

func (a *appVar) Alarm(title, body string) {
	title = Deploy.Name() + ` ` + title
	a.Mailer().Send(&mailer.Message{Receivers: a.Keeper(), Title: title, Body: body})
}

func (a *appVar) Keeper() []string {
	Load()
	return a.conf.Keeper
}
