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

func Get(env string) *Conf {
	if theConf == nil {
		if content, err := ioutil.ReadFile(filepath.Join(release.Root(), `deploy.yml`)); err != nil {
			log.Panic(err)
		} else {
			envConfs := map[string]*Conf{}
			if err = yaml.Unmarshal(content, &envConfs); err != nil {
				log.Panic(err)
			}
			theConf = envConfs[env]
		}
	}
	return theConf
}

func HasService(env, svcName string) bool {
	_, ok := Get(env).Services[svcName]
	return ok
}

func GetService(env, svcName string) Service {
	svc, ok := Get(env).Services[svcName]
	if !ok {
		log.Panicf(`deploy.yml: services.%s: undefined.`, svcName)
	}
	return svc
}

func ServiceNames(env string) (names []string) {
	services := Get(env).Services
	for _, svcName := range []string{`app`, `tasks`, `web`, `logc`, `cron`, `godoc`} {
		if _, ok := services[svcName]; ok {
			names = append(names, svcName)
		}
	}
	return
}

func ImageNameOf(env, svcName string) string {
	svc := GetService(env, svcName)
	if svc.Image == `` {
		log.Panicf(`deploy.yml: %s.image: empty.`, svcName)
	}
	return svc.Image
}

var rePort = regexp.MustCompile(`^\d+$`)

func InstancesOf(env, svcName string) (instances []string) {
	svc := GetService(env, svcName)
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

func FirstContainerNameOf(env, svcName string) string {
	name := release.AppConf(env).DeployName() + `_` + svcName
	if instances := InstancesOf(env, svcName); len(instances) > 0 {
		name += `.` + instances[0]
	}
	return name
}
