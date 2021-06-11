package release

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type Deploy struct {
	AccessNodes     map[string]string
	Services        map[string]*Service
	VolumesToCreate []string `yaml:"volumesToCreate"`
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
		if content, err := ioutil.ReadFile(filepath.Join(Root(), `deploy.yml`)); err != nil {
			log.Panic(err)
		} else {
			deploys = map[string]*Deploy{}
			if err = yaml.Unmarshal(content, &deploys); err != nil {
				log.Panic(err)
			}
		}
	}
	theDeploy := deploys[env]
	if theDeploy == nil {
		log.Fatalf(`deploy.yml: %s: undefined.`, env)
	}
	for name, svc := range theDeploy.Services {
		svc.env = env
		svc.name = name
	}
	return theDeploy
}

func HasService(svcName, env string) bool {
	_, ok := GetDeploy(env).Services[svcName]
	return ok
}

func GetService(svcName, env string) *Service {
	svc, ok := GetDeploy(env).Services[svcName]
	if !ok {
		log.Fatalf(`deploy.yml: %s.services.%s: undefined.`, env, svcName)
	}
	return svc
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

func ServiceNameRegexp(svcName, env string) string {
	var names []string
	if svcName != `` {
		names = []string{svcName}
	} else {
		names = ServiceNames(env)
	}
	var regex string

	switch len(names) {
	case 0:
		log.Fatalf(`deploy.yml: %s: no services defined.`, env)
	case 1:
		regex = names[0]
	default:
		regex = fmt.Sprintf(`(%s)`, strings.Join(names, `|`))
	}

	return `^/` + EnvConfig(env).DeployName() + `-` + regex + `(\.\d+)?`
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
	tag := time.Now().In(EnvConfig(env).TimeLocation).Format(`20060102-150405`)
	log.Println(`time tag: `, color.MagentaString(tag))
	return tag
}

func (svc Service) FirstContainerName() string {
	name := EnvConfig(svc.env).DeployName() + `-` + svc.name
	if ports := svc.Ports; len(ports) > 0 {
		name += `.` + strconv.FormatInt(int64(ports[0]), 10)
	}
	return name
}
