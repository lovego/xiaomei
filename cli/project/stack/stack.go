package stack

import (
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

type Stack struct {
	Version  string
	Registry string             `yaml:"-"`
	Services map[string]Service `yaml:",omitempty"`
	Volumes  map[string]Volume  `yaml:",omitempty"`
	Networks map[string]Network `yaml:",omitempty"`
}
type Service map[string]interface{}
type Volume map[string]interface{}
type Network map[string]interface{}

var theStack *Stack

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
	return theStack, nil
}

func (s Stack) ImageName(svcName string) string {
	return path.Join(s.Registry, config.Name(), svcName)
}

func Services() (map[string]Service, error) {
	stack, err := getStack()
	if err != nil {
		return nil, err
	}
	return stack.Services, nil
}
