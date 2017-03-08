package images

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

type Image struct {
	svcName string
	imageDriver
}

type imageDriver interface {
	Prepare() error
	BuildDir() string
	Dockerfile() string
	RunPorts() []string
	RunFiles() []string
}

func (i Image) Build(imgName string) error {
	if err := i.Prepare(); err != nil {
		return err
	}
	config.Log(color.GreenString(`building ` + i.svcName + ` image.`))
	_, err := cmd.Run(cmd.O{Dir: i.BuildDir()}, `docker`, `build`,
		`--file=`+i.Dockerfile(), `--tag=`+imgName, `.`,
	)
	return err
}

func (i Image) Run(imgName string) error {
	if err := i.PrepareOrBuild(imgName); err != nil {
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
	args = append(args, imgName)
	_, err := cmd.Run(cmd.O{}, `docker`, args...)
	return err
}

func (i Image) PrepareOrBuild(imgName string) error {
	if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true}, `docker`, `image`, `inspect`, imgName) {
		return i.Prepare()
	} else {
		return i.Build(imgName)
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
