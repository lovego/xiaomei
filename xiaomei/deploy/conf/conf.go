package conf

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/xiaomei/xiaomei/release"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Services        map[string]*Service
	VolumesToCreate []string `yaml:"volumesToCreate"`
}

type Service struct {
	name, env        string
	Nodes            map[string]string
	Image, Ports     string
	Command, Options []string
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
			if theConf == nil {
				log.Fatalf(`deploy.yml: %s: undefined.`, env)
			}
			for name, svc := range theConf.Services {
				svc.name = name
				svc.env = env
			}
		}
	}
	return theConf
}

func HasService(svcName, env string) bool {
	_, ok := Get(env).Services[svcName]
	return ok
}

func GetService(svcName, env string) *Service {
	svc, ok := Get(env).Services[svcName]
	if !ok {
		log.Fatalf(`deploy.yml: %s.services.%s: undefined.`, env, svcName)
	}
	return svc
}

func ServiceNames(env string) (names []string) {
	services := Get(env).Services
	for _, svcName := range []string{`app`, `tasks`, `web`, `logc`, `godoc`} {
		if _, ok := services[svcName]; ok {
			names = append(names, svcName)
		}
	}
	return
}

func (svc Service) ImageName() string {
	if svc.Image == `` {
		log.Panicf(`deploy.yml: %s.image: empty.`, svc.name)
	}
	return svc.Image
}

func (svc Service) ImageNameWithTag(timeTag string) string {
	if svc.Image == `` {
		log.Panicf(`deploy.yml: %s.image: empty.`, svc.name)
	}
	if timeTag == `` {
		return svc.Image
	} else {
		return svc.Image + `:` + svc.env + `-` + timeTag
	}
}

func TimeTag(env string) string {
	tag := time.Now().In(release.AppConf(env).TimeLocation).Format(`060102-150405`)
	log.Println(`time tag: `, color.MagentaString(tag))
	return tag
}

var rePort = regexp.MustCompile(`^\d+$`)

func (svc Service) Instances() (instances []string) {
	if svc.Ports == `` {
		return
	}
	for _, instance := range strings.Split(svc.Ports, `,`) {
		instance = strings.TrimSpace(instance)
		if rePort.MatchString(instance) {
			instances = append(instances, instance)
		} else {
			log.Panicf(`deploy.yml: %s.instances: illegal format.`, svc.name)
		}
	}
	return
}

func (svc Service) FirstContainerName() string {
	name := release.AppConf(svc.env).DeployName() + `-` + svc.name
	if instances := svc.Instances(); len(instances) > 0 {
		name += `.` + instances[0]
	}
	return name
}
