package config

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/mailer"
)

var App AppConf

type AppConf struct {
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
	Domain       string `yaml:"domain"`
	Secret       string `yaml:"secret"`
	DockerImage  string `yaml:"dockerImage"`
	StartTimeout string `yaml:"startTimeout"`

	TimeZone TimeZoneConf    `yaml:"timeZone"`
	Mailer   MailerConf      `yaml:"mailer"`
	Keepers  []mailer.People `yaml:"keepers"`

	GitAddr   string `yaml:"gitAddr"`
	GitBranch string `yaml:"gitBranch"`
}

type TimeZoneConf struct {
	Name   string `yaml:"name"`
	Offset int    `yaml:"offset"`
}

type MailerConf struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Sender mailer.People
	Passwd string `yaml:"passwd"`
}

func (a *AppConf) Root() string {
	if a.root == `` {
		if root := detectRoot(); root != `` {
			a.root = root
		} else {
			panic(`app root not found.`)
		}
	}
	return a.root
}

func (a *AppConf) Name() string {
	Load()
	return a.conf.Name
}

func (a *AppConf) Env() string {
	Load()
	return a.conf.Env
}

func (a *AppConf) DeployName() string {
	return a.Name() + `_` + a.Env()
}

func (a *AppConf) Bin() string {
	return filepath.Join(a.Root(), a.Name())
}

func (a *AppConf) Domain() string {
	Load()
	return a.conf.Domain
}

func (a *AppConf) Secret() string {
	Load()
	return a.conf.Secret
}

func (a *AppConf) DockerImage() string {
	Load()
	return a.conf.DockerImage
}

func (a *AppConf) StartTimeout() time.Duration {
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

func (a *AppConf) TimeZone() *time.Location {
	if a.timeZone == nil {
		Load()
		a.timeZone = time.FixedZone(a.conf.TimeZone.Name, a.conf.TimeZone.Offset)
	}
	return a.timeZone
}

func (a *AppConf) Mailer() *mailer.Mailer {
	a.mailer.Lock()
	defer a.mailer.Unlock()
	if !a.mailer.setted {
		Load()
		m := a.conf.Mailer
		a.mailer.Mailer = mailer.New(m.Host, m.Port, m.Sender, m.Passwd)
		a.mailer.setted = true
	}
	return a.mailer.Mailer
}

func (a *AppConf) Alarm(title, body string) {
	title = a.DeployName() + ` ` + title
	a.Mailer().Send(&mailer.Message{Receivers: a.Keepers(), Title: title, Body: body})
}

func (a *AppConf) Keepers() []mailer.People {
	Load()
	return a.conf.Keepers
}

func (a *AppConf) GitAddr() string {
	Load()
	return a.conf.GitAddr
}

func (a *AppConf) GitBranch() string {
	Load()
	if a.conf.GitBranch != `` {
		return a.conf.GitBranch
	}
	a.conf.GitBranch, _ = cmd.Run(cmd.O{Output: true, Panic: true},
		`git`, `rev-parse`, `--abbrev-ref`, `HEAD`)
	return a.conf.GitBranch
}
