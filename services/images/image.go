package images

import (
	"log"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/release"
)

type Image struct {
	svcName string
	image   interface{}
}

func (i Image) build(env, tag string, pull bool) error {
	if err := i.prepare(); err != nil {
		return err
	}
	log.Println(color.GreenString(`building ` + i.svcName + ` image.`))
	args := []string{`build`}
	if pull {
		args = append(args, `--pull`)
	}
	imgName := release.GetService(i.svcName, env).ImageName(tag)
	args = append(args, `--file=`+i.dockerfile(), `--tag=`+imgName, `.`)
	_, err := cmd.Run(cmd.O{Dir: i.buildDir(), Print: true}, `docker`, args...)
	return err
}

func (i Image) push(env, tag string) error {
	log.Println(color.GreenString(`pushing ` + i.svcName + ` image.`))
	imgName := release.GetService(i.svcName, env).ImageName(tag)
	_, err := cmd.Run(cmd.O{Print: true}, `docker`, `push`, imgName)
	return err
}

func (i Image) list(env string) error {
	_, err := cmd.Run(cmd.O{}, `docker`, `images`,
		`-f`, `reference=`+release.GetService(i.svcName, env).ImageName(``))
	return err
}
