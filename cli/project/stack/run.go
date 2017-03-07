package stack

import (
	"errors"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Run(svcName string) error {
	image, ok := imagesMap[svcName]
	if !ok {
		return errors.New(`no image registered for ` + svcName)
	}
	return image.Run()
}

func (i Image) Run() error {
	if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true}, `docker`, `image`, `inspect`, i.Name()) {
		if err := i.Prepare(); err != nil {
			return err
		}
	} else {
		if err := i.Build(); err != nil {
			return err
		}
	}
	args := []string{
		`run`, `--name=` + config.Name() + `_` + i.svcName, `-it`, `--rm`, `--network=host`, `--no-healthcheck`,
	}
	for _, mapping := range i.RunMapping() {
		args = append(args, `-v`, mapping)
	}
	args = append(args, i.Name())
	_, err := cmd.Run(cmd.O{}, `docker`, args...)
	return err
}
