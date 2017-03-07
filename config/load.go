package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils"
	"gopkg.in/yaml.v2"
)

func (c *Conf) Load() {
	c.Lock()
	defer c.Unlock()
	if c.data == nil {
		c.data = &conf{}
		Parse(c.data)
	}
}

func Parse(p interface{}) {
	loadConfig(p, `config/config.yml`)
	loadConfig(p, `config/envs/`+Env()+`.yml`)
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
