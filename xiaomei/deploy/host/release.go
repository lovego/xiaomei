package host

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
	"gopkg.in/yaml.v2"
)

const ConfigFile = `release.yml`

type releaseConf struct {
	Services        map[string]service
	VolumesToCreate []string
}

type service struct {
	Image, Ports string
	Volumes      []string
}

var theRelease *releaseConf

func getRelease() *releaseConf {
	if theRelease == nil {
		content, err := ioutil.ReadFile(filepath.Join(release.Root(), ConfigFile))
		if err != nil {
			panic(err)
		}
		conf := &releaseConf{}
		if err := yaml.Unmarshal(content, conf); err != nil {
			panic(err)
		}
		theRelease = conf
	}
	return theRelease
}

func getService(svcName string) service {
	svc, ok := getRelease().Services[svcName]
	if !ok {
		panic(fmt.Sprintf(`release.yml: services.%s: undefined.`, svcName))
	}
	return svc
}
