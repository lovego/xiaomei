package conf

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name string
	Envs map[string]*Conf
}

func Get(path string) *Config {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	config := &Config{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		log.Fatalf("parse %s: %v", path, err)
	}
	for env, conf := range config.Envs {
		conf.Name = config.Name
		conf.Env = env
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
