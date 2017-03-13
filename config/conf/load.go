package conf

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils"
	"gopkg.in/yaml.v2"
)

func New(root, env string) *Conf {
	conf := &Conf{root: root, env: env}
	Parse(&conf.data, root, env)
	return conf
}

func Parse(p interface{}, root, env string) {
	loadFile(p, filepath.Join(root, `config/config.yml`))
	loadFile(p, filepath.Join(root, `config/envs/`+env+`.yml`))
	if os.Getenv(`debugConf`) != `` {
		utils.PrintJson(p)
	}
}

func loadFile(p interface{}, file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(content, p)
	if err != nil {
		panic(err)
	}
}
