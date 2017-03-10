package release

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

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

func GetStack() *Stack {
	if theStack == nil {
		content, err := ioutil.ReadFile(filepath.Join(Root(), `stack.yml`))
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

func ImageNameOf(svcName string) string {
	service := GetStack().Services[svcName]
	if service == nil {
		panic(fmt.Sprintf(`stack.yml: services.%s: undefined.`, svcName))
	}
	image := service[`image`]
	if image == nil {
		panic(fmt.Sprintf(`stack.yml: services.%s.image: undefined.`, svcName))
	}
	if str, ok := image.(string); ok && str != `` {
		return str
	} else {
		panic(fmt.Sprintf(`stack.yml: services.%s.image: should be a string.`, svcName))
	}
}
