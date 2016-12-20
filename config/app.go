package config

import (
	"sync"
	"time"

	"github.com/bughou-go/xiaomei/utils/mailer"
)

var App appVar

type appVar struct {
	root string
	conf appConf
}

type appConf struct {
	Name         string `yaml:"name"`
	Env          string `yaml:"env"`
	Port         string `yaml:"port"`
	Domain       string `yaml:"domain"`
	Secret       string `yaml:"secret"`
	StartTimeout uint16 `yaml:"startTimeout"`

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

func (a *appVar) StartTimeout() uint16 {
	Load()
	return a.conf.StartTimeout
}

var timeZone struct {
	sync.Mutex
	*time.Location
}

func (a *appVar) TimeZone() *time.Location {
	timeZone.Lock()
	defer timeZone.Unlock()
	if timeZone.Location == nil {
		Load()
		timeZone.Location = time.FixedZone(a.conf.TimeZone.Name, a.conf.TimeZone.Offset)
	}
	return timeZone.Location
}

var _mailer struct {
	sync.Mutex
	setted bool
	*mailer.Mailer
}

func (a *appVar) Mailer() *mailer.Mailer {
	_mailer.Lock()
	defer _mailer.Unlock()
	if !_mailer.setted {
		Load()
		m := a.conf.Mailer
		if m.Host != `` && m.Port != `` && m.Sender != `` {
			_mailer.Mailer = mailer.New(m.Host, m.Port, m.Sender, m.Passwd)
		}
		_mailer.setted = true
	}
	return _mailer.Mailer
}

func (a *appVar) Alarm(title, body string) {
	title = Deploy.Name() + ` ` + title
	a.Mailer().Send(&mailer.Message{Receivers: a.Keeper(), Title: title, Body: body})
}

func (a *appVar) Keeper() []string {
	Load()
	return a.conf.Keeper
}
