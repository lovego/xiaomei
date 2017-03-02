package stack

import (
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

func ImageName(svcName string) (string, error) {
	if _, err := getStack(); err != nil {
		return ``, err
	} else {
		return path.Join(theRegistry, config.Name(), svcName), nil
	}
}

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

func getStack() (*Stack, error) {
	if theStack != nil {
		return theStack, nil
	}
	content, err := ioutil.ReadFile(filepath.Join(config.Root(), `../stack.yml`))
	if err != nil {
		return nil, err
	}
	stack := &Stack{Services: make(map[string]Service)}
	if err := yaml.Unmarshal(content, stack); err != nil {
		return nil, err
	}
	theStack = stack
	theRegistry = stack.Registry
	stack.Registry = ``
	return theStack, nil
}

func eachServiceDo(work func(svcName string) error) error {
	stack, err := getStack()
	if err != nil {
		return err
	}
	for svcName := range stack.Services {
		if svcName != `` {
			if err := work(svcName); err != nil {
				return err
			}
		}
	}
	return nil
}
