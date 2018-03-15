package conf

import (
	"io/ioutil"
	"log"

	"github.com/lovego/strmap"
	"gopkg.in/yaml.v2"
)

func Data(path string) strmap.StrMap {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	config := map[string]interface{}{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		log.Fatalf("parse %s: %v", path, err)
	}
	return strmap.StrMap(config)
}
