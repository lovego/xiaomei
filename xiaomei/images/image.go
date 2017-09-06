package images

import (
	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/registry"
)

type Image struct {
	svcName string
	image   interface{}
}

func (i Image) Build(env string, pull bool) error {
	if err := i.prepare(); err != nil {
		return err
	}
	utils.Log(color.GreenString(`building ` + i.svcName + ` image.`))
	args := []string{`build`}
	if pull {
		args = append(args, `--pull`)
	}
	args = append(args, `--file=`+i.dockerfile(), `--tag=`+conf.GetService(env, i.svcName).ImageNameAndTag(), `.`)
	_, err := cmd.Run(cmd.O{Dir: i.buildDir(), Print: true}, `docker`, args...)
	return err
}

func (i Image) Push(env string) error {
	utils.Log(color.GreenString(`pushing ` + i.svcName + ` image.`))
	_, err := cmd.Run(cmd.O{Print: true}, `docker`, `push`, conf.GetService(env, i.svcName).ImageNameAndTag())
	return err
}

// TODO: https or http check.
// TODO: https://registry.hub.docker.com/v2/
func (i Image) NameWithDigestInRegistry(env string) string {
	imgName := conf.GetService(env, i.svcName).ImageName()
	digest := registry.Digest(imgName, env)
	return imgName + `@` + digest
}
