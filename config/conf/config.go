package conf

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name string
	Envs map[string]*Conf
}

func Get(root string) *Config {
	path := filepath.Join(root, `config/config.yml`)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	config := &Config{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		log.Fatalf("parse config/config.yml: %v", err)
	}
	for _, conf := range config.Envs {
		conf.Name = config.Name
		if conf.TimeZone.Name != `` {
			conf.TimeLocation = time.FixedZone(conf.TimeZone.Name, conf.TimeZone.Offset)
		}
	}
	return config
}

func (config *Config) Get(env string) *Conf {
	conf := config.Envs[env]
	if conf == nil {
		log.Fatalf("config/config.yml envs.%s: not defined.", env)
	}
	return conf
}
