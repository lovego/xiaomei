package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/bughou-go/xiaomei/utils"
	"gopkg.in/yaml.v2"
)

var loader struct {
	mutex   sync.Mutex
	App     *appConf     `yaml:"app"`
	DB      *dbConf      `yaml:"db"`
	Deploy  *deployConf  `yaml:"deploy"`
	Servers *serversConf `yaml:"servers"`
	Godoc   *godocConf   `yaml:"godoc"`
}

func Load() {
	loader.mutex.Lock()
	defer loader.mutex.Unlock()
	if loader.App == nil {
		loader.App = &App.conf
		loader.DB = &DB.conf
		loader.Deploy = &Deploy.conf
		loader.Servers = &Servers.conf
		loader.Godoc = &Godoc.conf
		Parse(&loader)
	}
}

func Parse(p interface{}) {
	loadConfig(p, `config/config.yml`)
	loadConfig(p, envConfigPath())
	if Debug(`config`) {
		utils.PrintJson(p)
	}
}

func loadConfig(p interface{}, path string) {
	content, err := ioutil.ReadFile(filepath.Join(App.Root(), path))
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
