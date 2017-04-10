package host

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
)

// TODO: keep container history, wait until healthy
const deployScriptTmpl = `set -e
deploy() {
  name={{.Name}}{{ if .Ports }}.$1{{ end }}
  docker stop $name >/dev/null 2>&1 && docker rm $name
  docker run --name=$name {{ if .Ports }}-e {{.PortEnv}}=$1{{ end }} \
	{{ range .Envs }} -e {{ . }}{{ end }} \
  {{ range .Volumes}} -v {{ . }}{{ end }} \
  -d --network=host --restart=always \
	{{.Image}}
}
{{ range .VolumesToCreate -}}
docker volume create {{ . }}
{{ end }}
{{ if .Ports -}}
for port in {{ .Ports }}; do deploy $port; done
{{ else -}}
deploy
{{ end }}
exit
`

func getDeployScript(svcName string) (string, error) {
	tmpl := template.Must(template.New(``).Parse(deployScriptTmpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, getDeployConfig(svcName)); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

type deployConf struct {
	Name, Image, PortEnv, Ports    string
	Envs, VolumesToCreate, Volumes []string
}

func getDeployConfig(svcName string) deployConf {
	conf := deployConf{
		Name:            release.Name() + `_` + svcName,
		Image:           imageNameWithSha256Of(svcName),
		PortEnv:         portEnvName(svcName),
		Envs:            images.Get(svcName).EnvsForDeploy(),
		VolumesToCreate: getRelease().VolumesToCreate,
		Volumes:         getService(svcName).Volumes,
	}
	if conf.PortEnv != `` {
		conf.Ports = strings.Join(portsOf(svcName), ` `)
	}
	return conf
}

// TODO: https or http check.
// TODO: https://registry.hub.docker.com/v2/
func imageNameWithSha256Of(svcName string) string {
	imgName := Driver.ImageNameOf(svcName)
	uri, err := url.Parse(`http://` + imgName + `/manifests/latest`)
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
	return imgName + `@` + digest
}
