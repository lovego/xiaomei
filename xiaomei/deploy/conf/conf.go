package conf

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lovego/xiaomei/xiaomei/release"
	"gopkg.in/yaml.v2"
)

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
		theConf = &Conf{}
		if content, err := ioutil.ReadFile(filepath.Join(release.Root(), `deploy.yml`)); err != nil {
			log.Panic(err)
		} else {
			if err = yaml.Unmarshal(content, theConf); err != nil {
				log.Panic(err)
			}
		}
	}
	return theConf
}

func HasService(svcName string) bool {
	_, ok := Get().Services[svcName]
	return ok
}

func GetService(svcName string) Service {
	svc, ok := Get().Services[svcName]
	if !ok {
		log.Panicf(`deploy.yml: services.%s: undefined.`, svcName)
	}
	return svc
}

func ServiceNames() (names []string) {
	services := Get().Services
	for _, svcName := range []string{`app`, `tasks`, `web`, `logc`, `cron`, `godoc`} {
		if _, ok := services[svcName]; ok {
			names = append(names, svcName)
		}
	}
	return
}

func ImageNameOf(svcName string) string {
	svc := GetService(svcName)
	if svc.Image == `` {
		log.Panicf(`deploy.yml: %s.image: empty.`, svcName)
	}
	return svc.Image
}

var rePort = regexp.MustCompile(`^\d+$`)

func InstancesOf(svcName string) (instances []string) {
	svc := GetService(svcName)
	if svc.Ports == `` {
		return
	}
	for _, instance := range strings.Split(svc.Ports, `,`) {
		instance = strings.TrimSpace(instance)
		if rePort.MatchString(instance) {
			instances = append(instances, instance)
		} else {
			log.Panicf(`deploy.yml: %s.instances: illegal format.`, svcName)
		}
	}
	return
}

func FirstContainerNameOf(svcName, env string) string {
	name := release.AppConf(env).DeployName() + `_` + svcName
	if instances := InstancesOf(svcName); len(instances) > 0 {
		name += `.` + instances[0]
	}
	return name
}

func CommandFor(svcName string) []string {
	return GetService(svcName).Command
}

func OptionsFor(svcName string) []string {
	return GetService(svcName).Options
}
