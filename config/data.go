package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
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

	DeployUser    string         `yaml:"deployUser"`
	DeployRoot    string         `yaml:"deployRoot"`
	DeployServers []ServerConfig `yaml:"deployServers"`
	DeployName    string         `yaml:"deployName"`
	DeployPath    string         `yaml:"deployPath"`
	GitAddr       string         `yaml:"gitAddr"`
	GitBranch     string         `yaml:"gitBranch"`
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

func parseConfigData() *Config {
	data := &Config{}
	loadConfig(data, `config/config.yml`)
	// loadConfig(data, envConfigPath())
	utils.PrintJson(data)
	setupData(data)
	return data
}

func loadConfig(data *Config, p string) {
	content, err := ioutil.ReadFile(path.Join(Root, p))
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
	color.Yellow("%s %s\n", time.Now().Format(`2006-01-02 15:04:05 -0700`), configPath)

	return configPath
}

func setupData(data *Config) {
	data.DeployName = data.AppName + `_` + data.Env
	data.DeployPath = path.Join(data.DeployRoot, data.DeployName)
	if data.GitBranch != `` {
		return
	}
	branch, err := cmd.Run(cmd.O{Output: true, Panic: true},
		`git`, `rev-parse`, `--abbrev-ref`, `HEAD`)
	if err != nil {
		panic(err)
	}
	data.GitBranch = branch
}
