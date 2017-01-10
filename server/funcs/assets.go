package funcs

import (
	"html/template"
	"io/ioutil"
	"os"
	"path"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

var assets map[string]string

func init() {
	data, err := ioutil.ReadFile(path.Join(config.App.Root(), `config/assets.yml`))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &assets)
	if err != nil {
		panic(err)
	}
}

func AssetFunc(dev bool) func(string) string {
	return func(src string) string {
		if dev {
			return src + `?` + modificationTime(src)
		}
		if mt, ok := assets[src]; ok {
			return src + `?` + mt
		} else {
			return src
		}
	}
}

func modificationTime(src string) string {
	info, err := os.Stat(path.Join(config.App.Root(), `public`, src))
	if err != nil {
		panic(err)
	}
	return info.ModTime().Format(`060102150405`)
}

func HtmlSafe(text string) template.HTML {
	return template.HTML(text)
}
