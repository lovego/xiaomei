package funcs

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"github.com/bughou-go/xiaomei/config"
)

var assets map[string]string
var root = config.Root

func init() {
	data, err := ioutil.ReadFile(path.Join(root, `config/assets.json`))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &assets)
	if err != nil {
		panic(err)
	}
}

func AddModificationTimeFunc(cache bool) func(string) string {
	return func(src string) string {
		if cache {
			mt, ok := assets[src[7:]]
			if !ok {
				return src
			}
			return src + `?` + mt
		}
		return src + `?` + getModificationTime(src)
	}
}

func getModificationTime(src string) string {
	info, err := os.Stat(path.Join(root, `public/reports`, src[7:])) // remove "/static" prefix
	if err != nil {
		panic(err)
	}
	return info.ModTime().Format(`060102150405`)
}

func HtmlSafe(text string) template.HTML {
	return template.HTML(text)
}
