package swarm

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/lovego/xiaomei/xiaomei/release"
	"gopkg.in/yaml.v2"
)

const ConfigFile = `stack.yml`

type stack struct {
	Version  string
	Services map[string]service `yaml:",omitempty"`
	Volumes  map[string]volume  `yaml:",omitempty"`
	Networks map[string]network `yaml:",omitempty"`
}
type service map[string]interface{}
type volume map[string]interface{}
type network map[string]interface{}

var theStack *stack

func getStack() stack {
	if theStack == nil {
		content, err := ioutil.ReadFile(filepath.Join(release.Root(), ConfigFile))
		if err != nil {
			panic(err)
		}
		stack := &stack{Services: make(map[string]service)}
		if err := yaml.Unmarshal(content, stack); err != nil {
			panic(err)
		}
		theStack = stack
	}
	return *theStack
}

func getService(svcName string) service {
	svc := getStack().Services[svcName]
	if svc == nil {
		panic(fmt.Sprintf(`stack.yml: services.%s: undefined.`, svcName))
	}
	return svc
}
