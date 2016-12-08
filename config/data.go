package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/bughou-go/xiaomei/utils"
	"gopkg.in/yaml.v2"
)

type dataConfig struct {
	AppName string `yaml:"appName"`
	AppPort string `yaml:"appPort"`
	Env     string `yaml:"env"`
	Domain  string `yaml:"domain"`

	// for deploy
	DeployUser    string         `yaml:"deployUser"`
	DeployRoot    string         `yaml:"deployRoot"`
	DeployServers []ServerConfig `yaml:"deployServers"`
	GitAddr       string         `yaml:"gitAddr"`
	GitBranch     string         `yaml:"gitBranch"`

	// misc
	TimeZoneName   string       `yaml:"timeZoneName"`
	TimeZoneOffset int          `yaml:"timeZoneOffset"`
	Mailer         MailerConfig `yaml:"mailer"`
	AlarmReceivers []string     `yaml:"alarmReceivers"`

	// for db shell
	Mysql map[string]string `yaml:"mysql"`
	Redis map[string]string `yaml:"redis"`
}

type ServerConfig struct {
	Addr       string `yaml:"addr"`
	Tasks      string `yaml:"tasks"`
	AppAddr    string `yaml:"appAddr"`
	AppStartOn string `yaml:"appStartOn"`
	Misc       string `yaml:"misc"`
}

type MailerConfig struct {
	Host, Port, Sender, Passwd string
}

var dataInner struct {
	sync.Mutex
	*dataConfig
}

func data() *dataConfig {
	dataInner.Lock()
	defer dataInner.Unlock()
	if dataInner.dataConfig == nil {
		dataInner.dataConfig = &dataConfig{}
		Parse(dataInner.dataConfig)
	}
	return dataInner.dataConfig
}

func Parse(p interface{}) {
	loadConfig(p, `config/config.yml`)
	loadConfig(p, envConfigPath())
	if Debug(`config`) {
		utils.PrintJson(p)
	}
}

func loadConfig(p interface{}, path string) {
	content, err := ioutil.ReadFile(filepath.Join(Root(), path))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(content, p)
	if err != nil {
		panic(err)
	}
}

func envConfigPath() string {
	env := os.Getenv(`GOENV`)
	if env != `` {
		env = `envs/` + env
	} else {
		program := os.Args[0]
		if program[len(program)-5:] == `.test` {
			env = `envs/test`
		} else {
			env = `env`
		}
	}
	configPath := `config/` + env + `.yml`

	return configPath
}
