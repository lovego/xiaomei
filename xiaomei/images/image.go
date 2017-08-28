package images

import (
	"net/http"
	"net/url"

	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
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
	args = append(args, `--file=`+i.dockerfile(), `--tag=`+conf.ImageNameOf(env, i.svcName)+`:`+env, `.`)
	_, err := cmd.Run(cmd.O{Dir: i.buildDir(), Print: true}, `docker`, args...)
	return err
}

func (i Image) Push(env string) error {
	utils.Log(color.GreenString(`pushing ` + i.svcName + ` image.`))
	_, err := cmd.Run(cmd.O{Print: true}, `docker`, `push`, conf.ImageNameOf(env, i.svcName)+`:`+env)
	return err
}

// TODO: https or http check.
// TODO: https://registry.hub.docker.com/v2/
func (i Image) NameWithDigestInRegistry(env string) string {
	imgName := conf.ImageNameOf(env, i.svcName)
	uri, err := url.Parse(`http://` + imgName + `/manifests/` + env)
	if err != nil {
		panic(err)
	}
	uri.Path = `/v2` + uri.Path
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set(`Accept`, `application/vnd.docker.distribution.manifest.v2+json`)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	digest := resp.Header.Get(`Docker-Content-Digest`)
	if digest == `` {
		panic(`get image digest faild for: ` + imgName + `:` + env)
	}
	return imgName + `@` + digest
}
