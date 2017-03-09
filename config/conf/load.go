package conf

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/utils"
	"gopkg.in/yaml.v2"
)

func New(root, env string) Conf {
	conf := Conf{root: root, env: env}
	loadFile(&conf.data, filepath.Join(conf.root, `config/config.yml`))
	loadFile(&conf.data, filepath.Join(conf.root, `config/envs/`+env+`.yml`))
	if os.Getenv(`debugConf`) != `` {
		utils.PrintJson(&conf.data)
	}
	return conf
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
