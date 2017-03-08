package stack

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

type Stack struct {
	Version  string
	Services map[string]Service `yaml:",omitempty"`
	Volumes  map[string]Volume  `yaml:",omitempty"`
	Networks map[string]Network `yaml:",omitempty"`
}
type Service map[string]interface{}
type Volume map[string]interface{}
type Network map[string]interface{}

var theStack *Stack

func getStack() *Stack {
	if theStack == nil {
		content, err := ioutil.ReadFile(filepath.Join(config.Root(), `../stack.yml`))
		if err != nil {
			panic(err)
		}
		stack := &Stack{Services: make(map[string]Service)}
		if err := yaml.Unmarshal(content, stack); err != nil {
			panic(err)
		}
		theStack = stack
	}
	return theStack
}

func eachServiceDo(work func(svcName, imgName string) error) error {
	for svcName := range getStack().Services {
		if svcName != `` {
			if imgName, err := serviceImageName(svcName); err != nil {
				return err
			} else if err := work(svcName, imgName); err != nil {
				return err
			}
		}
	}
	return nil
}

func serviceImageName(svcName string) (string, error) {
	service := getStack().Services[svcName]
	if service == nil {
		return ``, fmt.Errorf(`stack.yml: services.%s: undefined.`, svcName)
	}
	image := service[`image`]
	if image == nil {
		return ``, fmt.Errorf(`stack.yml: services.%s.image: undefined.`, svcName)
	}
	if str, ok := image.(string); ok && str != `` {
		return str, nil
	} else {
		return ``, fmt.Errorf(`stack.yml: services.%s.image: should be a string.`, svcName)
	}
}
