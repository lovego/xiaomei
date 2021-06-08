package images

import (
	"log"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services/images/app"
	"github.com/lovego/xiaomei/services/images/logc"
	"github.com/lovego/xiaomei/services/images/web"
)

var imagesMap = map[string]Image{
	`app`:  {svcName: `app`, image: app.Image{}},
	`web`:  {svcName: `web`, image: web.Image{}},
	`logc`: {svcName: `logc`, image: logc.Image{}},
}

type Image struct {
	svcName string
	image   interface{}
}

func Get(svcName string) Image {
	if img, ok := imagesMap[svcName]; !ok {
		panic(`no image for: ` + svcName)
	} else {
		return img
	}
}

type Build struct {
	Env, Tag         string
	GoBuildFlags     []string
	DockerBuildFlags []string
}

func (b Build) Run(svcName string) error {
	return imagesDo(svcName, b.Env, func(img Image) error {
		if err := img.prepare(b.GoBuildFlags); err != nil {
			return err
		}

		log.Println(color.GreenString(`building ` + img.svcName + ` image.`))

		_, err := cmd.Run(cmd.O{Dir: img.buildDir(), Print: true}, `docker`, b.args(img)...)
		return err
	})
}

func (b Build) args(img Image) []string {
	var result = []string{
		`build`,
		`--tag=` + release.GetService(img.svcName, b.Env).ImageName(b.Tag),
	}
	if len(b.DockerBuildFlags) > 0 {
		b.DockerBuildFlags = []string{`--pull`, `--file=` + img.dockerfile()}
	}
	result = append(result, b.DockerBuildFlags...)

	return append(result, `.`)
}

func Push(svcName, env, tag string) error {
	return imagesDo(svcName, env, func(img Image) error {
		log.Println(color.GreenString(`pushing ` + img.svcName + ` image.`))
		imgName := release.GetService(img.svcName, env).ImageName(tag)
		_, err := cmd.Run(cmd.O{Print: true}, `docker`, `push`, imgName)
		return err
	})
}

func List(svcName, env string) error {
	return imagesDo(svcName, env, func(img Image) error {
		if svcName == `` {
			color.Green(img.svcName + `:`)
		}
		_, err := cmd.Run(cmd.O{}, `docker`, `images`,
			`-f`, `reference=`+release.GetService(img.svcName, env).ImageName(``))
		return err
	})
}

func imagesDo(svcName, env string, work func(Image) error) error {
	if svcName == `` {
		for _, svcName := range release.ServiceNames(env) {
			if err := work(Get(svcName)); err != nil {
				return err
			}
		}
		return nil
	} else {
		return work(Get(svcName))
	}
}
