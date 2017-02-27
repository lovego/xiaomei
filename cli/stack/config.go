package stack

import (
	"io/ioutil"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

type Stack struct {
	Services map[string]Service
}

var stack = load()

func load() Stack {
	content, err := ioutil.ReadFile(filepath.Join(config.Root(), `../stack.yml`))
	if err != nil {
		panic(err)
	}
	stack := Stack{}
	if err = yaml.Unmarshal(content, &stack); err != nil {
		panic(err)
	}
	return stack
}
