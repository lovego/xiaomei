package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bughou-go/xiaomei/utils"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	App     *AppConf
	Cluster *ClusterConf
	Db      *DbConf
	Fmwk    *FmwkConf
	Godoc   *GodocConf
}

func Data() Conf {
	return Conf{&App, &Cluster, &DB, &Fmwk, &Godoc}
}

var loader struct {
	mutex   sync.Mutex
	App     *appConf     `yaml:"app"`
	DB      *dbConf      `yaml:"db"`
	Cluster *clusterConf `yaml:"cluster"`
	Godoc   *godocConf   `yaml:"godoc"`
}

func Load() {
	loader.mutex.Lock()
	defer loader.mutex.Unlock()
	if loader.App == nil {
		loader.App = &App.conf
		loader.DB = &DB.conf
		loader.Cluster = &Cluster.conf
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
	if env == `` {
		if program := os.Args[0]; strings.HasSuffix(program, `.test`) {
			env = `test`
		} else {
			env = `dev`
		}
	}
	configPath := `config/envs/` + env + `.yml`

	return configPath
}
