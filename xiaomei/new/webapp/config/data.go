package config

import (
	"sync"

	"github.com/bughou-go/xiaomei/config"
)

type Config struct {
	CustomKey1 string `yaml:"customKey1"`
	CustomKey2 string `yaml:"customKey2"`
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
