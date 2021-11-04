package release

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/config/config"
	"gopkg.in/yaml.v2"
)

type Deploy struct {
	AccessNodes     map[string]string   `yaml:"accessNodes"`
	Services        map[string]*Service `yaml:"services"`
	VolumesToCreate []string            `yaml:"volumesToCreate"`
}

type Service struct {
	env, name        string
	Nodes            map[string]string
	Image            string
	Ports            []uint16
	Command, Options []string
}

var deploys map[string]*Deploy

func GetDeploy(env string) *Deploy {
	if deploys == nil {
		if content, err := ioutil.ReadFile(configFile(env, `deploy.yml`)); err != nil {
			log.Panic(err)
		} else {
			deploys = map[string]*Deploy{}
			if err = yaml.Unmarshal(content, &deploys); err != nil {
				log.Panic(err)
			}
		}
	}
	environ := config.NewEnv(env)
	theDeploy := deploys[environ.Minor()]
	if theDeploy == nil {
		log.Fatalf(`%s: %s: undefined.`, configFile(env, `deploy.yml`), environ.Minor())
	}
	for name, svc := range theDeploy.Services {
		svc.env = env
		svc.name = name
	}
	return theDeploy
}

func HasService(env, svcName string) bool {
	_, ok := GetDeploy(env).Services[svcName]
	return ok
}

func GetService(env, svcName string) *Service {
	svc, ok := GetDeploy(env).Services[svcName]
	if !ok {
		log.Fatalf(`%s: %s.services.%s: undefined.`, configFile(env, `deploy.yml`), env, svcName)
	}
	return svc
}

func MultiPorts(env string, svcNames []string) bool {
	for _, svcName := range svcNames {
		if len(GetService(env, svcName).Ports) >= 2 {
			return true
		}
	}
	return false
}

func ServiceNames(env string) (names []string) {
	services := GetDeploy(env).Services
	for _, svcName := range []string{`app`, `web`, `logc`} {
		if _, ok := services[svcName]; ok {
			names = append(names, svcName)
		}
	}
	return
}

func ContainerNameRegexp(svcName, env string) string {
	var names []string
	if svcName != `` {
		names = []string{svcName}
	} else {
		names = ServiceNames(env)
	}
	var svcNamesRegexp string
	switch len(names) {
	case 0:
		log.Fatalf(`deploy.yml: %s: no services defined.`, env)
	case 1:
		svcNamesRegexp = names[0]
	default:
		svcNamesRegexp = fmt.Sprintf(`(%s)`, strings.Join(names, `|`))
	}

	return `^/` + regexp.QuoteMeta(Config(env).DeployName()) + `\.` + svcNamesRegexp + `(\.\d+)?$`
}

func (svc Service) ImageName(tag string) string {
	if svc.Image == `` {
		log.Panicf(`deploy.yml: %s.image: empty.`, svc.name)
	}
	if tag == `` {
		return svc.Image
	} else {
		return svc.Image + `:` + tag
	}
}

func TimeTag(env string) string {
	tag := time.Now().In(Config(env).TimeLocation).Format(`20060102-150405`)
	log.Println(`time tag: `, color.MagentaString(tag))
	return tag
}

func (svc Service) FirstContainerName() string {
	name := ServiceName(svc.env, svc.name)
	if ports := svc.Ports; len(ports) > 0 {
		name += `.` + strconv.FormatInt(int64(ports[0]), 10)
	}
	return name
}
