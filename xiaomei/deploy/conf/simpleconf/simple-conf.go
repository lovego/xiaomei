package simpleconf

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
	"gopkg.in/yaml.v2"
)

const File = `simple.yml`

type Conf struct {
	loaded          bool
	Services        map[string]Service
	VolumesToCreate []string
}

type Service struct {
	Image, Ports string
	Volumes      []string
}

var theConf *Conf

func Get() *Conf {
	if theConf == nil {
		content, err := ioutil.ReadFile(filepath.Join(release.Root(), File))
		if err != nil {
			panic(err)
		}
		conf := &Conf{}
		if err := yaml.Unmarshal(content, conf); err != nil {
			panic(err)
		}
		theConf = conf
	}
	return theConf
}

func GetService(svcName string) Service {
	svc, ok := Get().Services[svcName]
	if !ok {
		panic(fmt.Sprintf(`simple.yml: services.%s: undefined.`, svcName))
	}
	return svc
}

func ServiceNames() map[string]bool {
	m := make(map[string]bool)
	for svcName := range Get().Services {
		m[svcName] = true
	}
	return m
}

func ImageNameOf(svcName string) string {
	svc := GetService(svcName)
	if svc.Image == `` {
		panic(fmt.Sprintf(`release.yml: %s.image: empty.`, svcName))
	}
	return svc.Image
}
