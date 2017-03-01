package project

import (
	"fmt"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func GetService(name string) (Service, error) {
	stack, err := GetStack()
	if err != nil {
		return nil, err
	}
	if svc, ok := stack.Services[name]; ok {
		svc[`name`] = name
		return svc, nil
	} else {
		return nil, fmt.Errorf(`services.%s: undefined.`, name)
	}
}

func (s Service) GetImage() (string, error) {
	if v := s[`image`]; v != nil {
		if img, ok := v.(string); ok && img != `` {
			return img, nil
		} else if !ok {
			return ``, fmt.Errorf(`services.%s.image: must be a string. `, s[`name`])
		}
	}
	return ``, fmt.Errorf(`services.%s.image: undefined. `, s[`name`])
}

func (s Service) BuildImage() error {
	image, err := s.GetImage()
	if err != nil {
		return err
	}
	context, dockerfile, err := s.GetBuild()
	if err != nil {
		return err
	}
	context = filepath.Join(config.Root(), `..`, context)
	_, err = cmd.Run(cmd.O{Dir: context}, `docker`, `build`,
		`--file=`+dockerfile, `--tag=`+image, `.`)
	return err
}

func (s Service) GetBuild() (context, dockerfile string, err error) {
	switch build := s[`build`].(type) {
	case string:
		context = build
	case map[string]interface{}:
		for key, ifc := range build {
			if key != `context` && key != `dockerfile` {
				return ``, ``, fmt.Errorf(`services.%s.build.%s: unknown key.`, s[`name`], key)
			}
			value, ok := ifc.(string)
			if !ok {
				return ``, ``, fmt.Errorf(`services.%s.build.%s: should be a string.`, s[`name`], key)
			}
			switch key {
			case `context`:
				context = value
			case `dockerfile`:
				dockerfile = value
			}
		}
	default:
		return ``, ``, fmt.Errorf(`services.%s.build: should be a string or map[string]string.`, s[`name`])
	}
	return
}
