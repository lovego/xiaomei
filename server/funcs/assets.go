package funcs

import (
	"html/template"
	"io/ioutil"
	"path"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

var assets map[string]string

func init() {
	data, err := ioutil.ReadFile(path.Join(config.Root(), `config/assets.yml`))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &assets)
	if err != nil {
		panic(err)
	}
}

func AssetFunc() func(string) string {
	dev := config.DevMode()
	return func(src string) string {
		if dev {
			return src + `?` + time.Now().Format(`060102150405`)
		}
		if mt, ok := assets[src]; ok {
			return src + `?` + mt
		} else {
			return src
		}
	}
}

func HtmlSafe(text string) template.HTML {
	return template.HTML(text)
}
