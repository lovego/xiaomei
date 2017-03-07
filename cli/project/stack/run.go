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
	if err := i.PrepareOrBuild(); err != nil {
		return err
	}
	networkName := config.Name() + `_run`
	if err := ensureNetwork(networkName); err != nil {
		return err
	}
	args := []string{
		`run`, `-it`, `--rm`, `--no-healthcheck`,
		`--name=` + config.Name() + `_` + i.svcName,
		`--network=` + networkName,
	}
	for _, port := range i.RunPorts() {
		args = append(args, `-p`, port)
	}
	for _, file := range i.RunFiles() {
		args = append(args, `-v`, file)
	}
	args = append(args, i.Name())
	_, err := cmd.Run(cmd.O{}, `docker`, args...)
	return err
}

func (i Image) PrepareOrBuild() error {
	if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true}, `docker`, `image`, `inspect`, i.Name()) {
		return i.Prepare()
	} else {
		return i.Build()
	}
}

func ensureNetwork(name string) error {
	if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true}, `docker`, `network`, `inspect`, name) {
		return nil
	}
	_, err := cmd.Run(cmd.O{}, `docker`, `network`, `create`,
		`--attachable`, `--driver=overlay`, name)
	return err
}
