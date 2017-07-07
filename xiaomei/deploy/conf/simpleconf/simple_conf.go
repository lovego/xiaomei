package simpleconf

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lovego/xiaomei/xiaomei/release"
	"gopkg.in/yaml.v2"
)

const File = `simple.yml`

type Conf struct {
	Services        map[string]Service
	VolumesToCreate []string `yaml:"volumesToCreate"`
}

type Service struct {
	Nodes              map[string]string
	Name, Ports, Image string
	Command, Options   []string
}

var theConf *Conf

func Get() *Conf {
	if theConf == nil {
		envConfs := map[string]*Conf{}
		if content, err := ioutil.ReadFile(filepath.Join(release.Root(), File)); err != nil {
			panic(err)
		} else {
			if err = yaml.Unmarshal(content, &envConfs); err != nil {
				panic(err)
			}
		}
		theConf = envConfs[release.Env()]
	}
	return theConf
}

func HasService(svcName string) bool {
	_, ok := Get().Services[svcName]
	return ok
}

func GetService(svcName string) Service {
	if svc, ok := Get().Services[svcName]; ok {
		return svc
	}
	panic(fmt.Sprintf(`simple.yml: services.%s: undefined.`, svcName))
}

func ServiceNames() (names []string) {
	services := Get().Services
	for _, svcName := range []string{`app`, `tasks`, `web`, `logc`, `godoc`} {
		if _, ok := services[svcName]; ok {
			names = append(names, svcName)
		}
	}
	return
}

func ImageNameOf(svcName string) string {
	svc := GetService(svcName)
	if svc.Image == `` {
		panic(fmt.Sprintf(`%s: %s.image: empty.`, File, svcName))
	}
	return svc.Image
}

var rePort = regexp.MustCompile(`^\d+$`)

func PortsOf(svcName string) (ports []string) {
	svc := GetService(svcName)
	if svc.Ports == `` {
		return
	}
	for _, port := range strings.Split(svc.Ports, `,`) {
		port = strings.TrimSpace(port)
		if rePort.MatchString(port) {
			ports = append(ports, port)
		} else {
			panic(fmt.Sprintf(`%s: %s.ports: illegal format.`, File, svcName))
		}
	}
	return
}

func ContainerNameOf(svcName string) string {
	name := release.Name() + `_` + svcName
	if ports := PortsOf(svcName); len(ports) > 0 {
		name += `.` + ports[0]
	}
	return name
}

func CommandFor(svcName string) []string {
	return GetService(svcName).Command
}

func OptionsFor(svcName string) []string {
	return GetService(svcName).Options
}
