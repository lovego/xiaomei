package config

import (
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AppName string `yaml:"appName"`
	AppPort string `yaml:"appPort"`
	Env     string `yaml:"env"`
	Domain  string `yaml:"domain"`

	TimeZoneName   string       `yaml:"timeZoneName"`
	TimeZoneOffset int          `yaml:"timeZoneOffset"`
	Mailer         MailerConfig `yaml:"mailer"`
	AlarmReceivers []string     `yaml:"alarmReceivers"`

	// for deploy
	DeployUser    string         `yaml:"deployUser"`
	DeployRoot    string         `yaml:"deployRoot"`
	DeployServers []ServerConfig `yaml:"deployServers"`
	DeployName    string         `yaml:"deployName"`
	DeployPath    string         `yaml:"deployPath"`
	GitAddr       string         `yaml:"gitAddr"`
	GitBranch     string         `yaml:"gitBranch"`

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

var data struct {
	sync.Mutex
	*Config
}

func Data() Config {
	data.Lock()
	defer data.Unlock()
	if data.Config == nil {
		data.Config = &Config{}
		Parse(data.Config)
	}
	return *data.Config
}

func Parse(data interface{}) {
	loadConfig(data, `config/config.yml`)
	loadConfig(data, envConfigPath())
	if d, ok := data.(*Config); ok {
		setupData(d)
	} else {
	}
	if Debug(`config`) {
		utils.PrintJson(data)
	}
}

func loadConfig(data interface{}, p string) {
	content, err := ioutil.ReadFile(path.Join(Root(), p))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(content, data)
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

func setupData(data *Config) {
	data.DeployName = data.AppName + `_` + data.Env
	data.DeployPath = path.Join(data.DeployRoot, data.DeployName)
	if data.GitBranch == `` {
		branch, err := cmd.Run(cmd.O{Output: true, Panic: true},
			`git`, `rev-parse`, `--abbrev-ref`, `HEAD`)
		if err != nil {
			panic(err)
		}
		data.GitBranch = branch
	}
}
