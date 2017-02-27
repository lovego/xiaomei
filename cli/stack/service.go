package stack

import (
	"fmt"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func GetService(name string) (Service, error) {
	stack, err := GetStack()
	if err != nil {
		return Service{}, err
	}
	if svc, ok := stack.Services[name]; ok {
		svc.name = name
		return svc, nil
	} else {
		return Service{}, fmt.Errorf(`services.%s: undefined.`, name)
	}
}

type Service struct {
	name  string
	Image string
	Build interface{}
}

func (s Service) GetImage() (string, error) {
	if img := s.Image; img != `` {
		return img, nil
	} else {
		return ``, fmt.Errorf(`services.%s.image: undefined. `, s.name)
	}
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
		`--dockerfile=`+dockerfile, `--tag=`+image, `.`)
	return err
}

func (s Service) GetBuild() (context, dockerfile string, err error) {
	switch build := s.Build.(type) {
	case string:
		context = build
	case map[string]interface{}:
		for key, ifc := range build {
			if key != `context` && key != `dockerfile` {
				return ``, ``, fmt.Errorf(`services.%s.build.%s: unknown key.`, s.name, key)
			}
			value, ok := ifc.(string)
			if !ok {
				return ``, ``, fmt.Errorf(`services.%s.build.%s: should be a string.`, s.name, key)
			}
			switch key {
			case `context`:
				context = value
			case `dockerfile`:
				dockerfile = value
			}
		}
	default:
		return ``, ``, fmt.Errorf(`services.%s.build: should be a string or map[string]string.`, s.name)
	}
	return
}
