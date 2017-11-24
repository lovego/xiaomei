package images

import (
	"log"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
)

type Image struct {
	svcName string
	image   interface{}
}

func (i Image) build(env, timeTag string, pull bool) error {
	if err := i.prepare(); err != nil {
		return err
	}
	log.Println(color.GreenString(`building ` + i.svcName + ` image.`))
	args := []string{`build`}
	if pull {
		args = append(args, `--pull`)
	}
	imgName := conf.GetService(i.svcName, env).ImageNameWithTag(timeTag)
	args = append(args, `--file=`+i.dockerfile(), `--tag=`+imgName, `.`)
	_, err := cmd.Run(cmd.O{Dir: i.buildDir(), Print: true}, `docker`, args...)
	return err
}

func (i Image) push(env, timeTag string) error {
	log.Println(color.GreenString(`pushing ` + i.svcName + ` image.`))
	imgName := conf.GetService(i.svcName, env).ImageNameWithTag(timeTag)
	_, err := cmd.Run(cmd.O{Print: true}, `docker`, `push`, imgName)
	return err
}
