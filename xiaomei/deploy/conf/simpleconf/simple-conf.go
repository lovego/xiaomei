package simpleconf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
	"github.com/lovego/xiaomei/xiaomei/release"
	"gopkg.in/yaml.v2"
)

const File = `simple.yml`

type Conf struct {
	Services        map[string]Service
	VolumesToCreate []string `yaml:"volumesToCreate"`
}

type Service struct {
	Image, Ports     string
	Command, Volumes []string
}

var theConf *Conf

func Get() *Conf {
	if theConf == nil {
		conf := &Conf{}
		loadFile(conf, filepath.Join(release.Root(), File))
		if envConf := filepath.Join(release.Root(), `simple/`+release.Env()+`.yml`); fs.Exist(envConf) {
			loadFile(conf, envConf)
		}
		if os.Getenv(`debugSimpleConf`) != `` {
			utils.PrintJson(conf)
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
		panic(fmt.Sprintf(`%s: %s.image: empty.`, File, svcName))
	}
	return svc.Image
}

var rePort = regexp.MustCompile(`^\d+$`)
var rePorts = regexp.MustCompile(`^(\d+)-(\d+)$`)

func PortsOf(svcName string) (ports []string) {
	svc := GetService(svcName)
	if svc.Ports == `` {
		return
	}
	if rePort.MatchString(svc.Ports) {
		ports = append(ports, svc.Ports)
	} else if m := rePorts.FindStringSubmatch(svc.Ports); len(m) == 3 {
		start, err := strconv.Atoi(m[1])
		if err != nil {
			panic(err)
		}
		end, err := strconv.Atoi(m[2])
		if err != nil {
			panic(err)
		}
		for ; start <= end; start++ {
			ports = append(ports, strconv.Itoa(start))
		}
	} else {
		panic(fmt.Sprintf(`%s: %s.ports: illegal format.`, File, svcName))
	}
	return
}

func CommandFor(svcName string) []string {
	return GetService(svcName).Command
}

func VolumesFor(svcName string) []string {
	return GetService(svcName).Volumes
}

func loadFile(p interface{}, file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(content, p)
	if err != nil {
		panic(err)
	}
}
