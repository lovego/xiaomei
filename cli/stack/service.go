package stack

import (
	"errors"
)

type Service struct {
	Image string
	Build interface{}
}

func ServiceImage(name string) (string, error) {
	if svc, ok := stack.Services[name]; ok {
		if img := svc.Image; img != `` {
			return img, nil
		} else {
			return ``, errors.New(`empty img for service: ` + name)
		}
	} else {
		return ``, errors.New(`no such service: ` + name)
	}
}
