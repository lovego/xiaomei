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

var reImageName = regexp.MustCompile(`^(.+):([\w.-]+)$`)

func ImageNameOf(svcName string) string {
	name := imageNameOf(svcName)
	if !reImageName.MatchString(name) {
		name += `:` + release.Env()
	}
	return name
}

func ImageNameAndTagOf(svcName string) (name, tag string) {
	name = imageNameOf(svcName)
	if m := reImageName.FindStringSubmatch(name); len(m) == 3 {
		return m[1], m[2]
	} else {
		return name, release.Env()
	}
}

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
		if content, err := ioutil.ReadFile(filepath.Join(release.Root(), `simple.yml`)); err != nil {
			log.Panic(err)
		} else {
			if err = yaml.Unmarshal(content, &envConfs); err != nil {
				log.Panic(err)
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
	svc, ok := Get().Services[svcName]
	if !ok {
		log.Panicf(`simple.yml: services.%s: undefined.`, svcName)
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

func imageNameOf(svcName string) string {
	svc := GetService(svcName)
	if svc.Image == `` {
		log.Panicf(`simple.yml: %s.image: empty.`, svcName)
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
			log.Panicf(`simple.yml: %s.instances: illegal format.`, svcName)
		}
	}
	return
}

func ContainerNameOf(svcName string) string {
	name := release.DeployName() + `_` + svcName
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
