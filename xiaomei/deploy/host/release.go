package host

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
	"gopkg.in/yaml.v2"
)

const ConfigFile = `release.yml`

type service struct {
	Image, Ports string
	Volumes      []string
}

var theRelease map[string]service

func getRelease() map[string]service {
	if theRelease == nil {
		content, err := ioutil.ReadFile(filepath.Join(release.Root(), ConfigFile))
		if err != nil {
			panic(err)
		}
		release := make(map[string]service)
		if err := yaml.Unmarshal(content, release); err != nil {
			panic(err)
		}
		theRelease = release
	}
	return theRelease
}

func getService(svcName string) service {
	svc, ok := getRelease()[svcName]
	if !ok {
		panic(fmt.Sprintf(`release.yml: %s: undefined.`, svcName))
	}
	return svc
}
