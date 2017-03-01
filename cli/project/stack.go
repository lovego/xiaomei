package project

import (
	"io/ioutil"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

type Stack struct {
	Version  string
	Services map[string]Service
	Volumes  map[string]Volume  `yaml:",omitempty"`
	Networks map[string]Network `yaml:",omitempty"`
}
type Service map[string]interface{}
type Volume map[string]interface{}
type Network map[string]interface{}

var theStack *Stack

func GetStack() (*Stack, error) {
	if theStack != nil {
		return theStack, nil
	}
	content, err := GetStackFileContent()
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

func GetStackFileContent() ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(config.Root(), `../stack.yml`))
}
