package stack

import (
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

func ImageName(svcName string) (string, error) {
	if stack, err := getStack(); err != nil {
		return ``, err
	} else {
		return stack.ImageName(svcName), nil
	}
}

func BuildImage(svcName string) error {
	config.Log(color.GreenString(`building ` + svcName + ` image.`))

	imageName, err := ImageName(svcName)
	if err != nil {
		return err
	}

	var dir, file string
	if svcName == `cron` {
		dir, file = `img-app`, `DockerfileCron`
	} else {
		dir, file = `img-`+svcName, `Dockerfile`
	}
	dir = filepath.Join(config.Root(), `..`, dir)

	_, err = cmd.Run(cmd.O{Dir: dir}, `docker`, `build`, `--file=`+file, `--tag=`+imageName, `.`)
	return err
}
