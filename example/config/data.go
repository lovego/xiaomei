package config

import (
	"sync"

	"github.com/bughou-go/xiaomei/config"
)

type Config struct {
	Mysql string `yaml:"mysql"`
	Redis string `yaml:"redis"`
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
		config.Parse(data.Config)
	}
	return *data.Config
}
