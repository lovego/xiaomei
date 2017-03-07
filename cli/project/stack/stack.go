package stack

import (
	"io/ioutil"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

type Stack struct {
	Version  string
	Registry string             `yaml:",omitempty"`
	Services map[string]Service `yaml:",omitempty"`
	Volumes  map[string]Volume  `yaml:",omitempty"`
	Networks map[string]Network `yaml:",omitempty"`
}
type Service map[string]interface{}
type Volume map[string]interface{}
type Network map[string]interface{}

var theStack *Stack
var theRegistry string

func getStack() *Stack {
	loadStack()
	return theStack
}

func getRegistry() string {
	loadStack()
	return theRegistry
}

func loadStack() {
	if theStack != nil {
		return
	}
	content, err := ioutil.ReadFile(filepath.Join(config.Root(), `../stack.yml`))
	if err != nil {
		panic(err)
	}
	stack := &Stack{Services: make(map[string]Service)}
	if err := yaml.Unmarshal(content, stack); err != nil {
		panic(err)
	}
	theStack = stack
	theRegistry = stack.Registry
	stack.Registry = ``
}

func eachServiceDo(work func(svcName string) error) error {
	for svcName := range getStack().Services {
		if svcName != `` {
			if err := work(svcName); err != nil {
				return err
			}
		}
	}
	return nil
}
