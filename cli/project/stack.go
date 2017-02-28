package project

import (
	"io/ioutil"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

type Stack struct {
	Services map[string]Service
}

var theStack *Stack

func GetStack() (*Stack, error) {
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
