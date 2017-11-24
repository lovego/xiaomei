package conf

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/lovego/utils/strmap"
	"gopkg.in/yaml.v2"
)

func Data(root, env string) strmap.StrMap {
	path := filepath.Join(root, `config/envs/`+env+`.yml`)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	config := map[string]interface{}{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		log.Fatalf("parse config/envs/%s.yml: %v", env, err)
	}
	return strmap.StrMap(config)
}
