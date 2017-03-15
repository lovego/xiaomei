package images

import (
	"encoding/json"
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/release"
)

func (i Image) Run(publish []string) error {
	if err := i.prepareForRun(); err != nil {
		return err
	}
	args := []string{`run`, `-it`, `--rm`,
		`--name=` + release.Name() + `_` + i.svcName, `--network=` + i.networkNameForRun(),
	}
	if pubs, err := i.publishForRun(publish); err == nil {
		for _, pub := range pubs {
			args = append(args, `-p`, pub)
		}
	} else {
		return err
	}
	for _, file := range i.FilesForRun() {
		args = append(args, `-v`, file)
	}
	for _, env := range i.EnvForRun() {
		args = append(args, `-e`, env)
	}
	args = append(args, release.ImageNameOf(i.svcName))
	if cmd := i.CmdForRun(); cmd != nil {
		args = append(args, cmd...)
	}
	_, err := cmd.Run(cmd.O{}, `docker`, args...)
	return err
}

func (i Image) prepareForRun() error {
	if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true},
		`docker`, `image`, `inspect`, release.ImageNameOf(i.svcName)) {
		if err := i.PrepareForBuild(); err != nil {
			return err
		}
	} else {
		if err := i.Build(); err != nil {
			return err
		}
	}
	return i.ensureNetworkForRun()
}

func (i Image) ensureNetworkForRun() error {
	name := i.networkNameForRun()
	if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true}, `docker`, `network`, `inspect`, name) {
		return nil
	}
	_, err := cmd.Run(cmd.O{}, `docker`, `network`, `create`,
		`--attachable`, `--driver=overlay`, name)
	return err
}

func (i Image) networkNameForRun() string {
	return release.Name() + `_run`
}

func (i Image) publishForRun(publish []string) ([]string, error) {
	if len(publish) == 0 {
		if publish = release.PortsOf(i.svcName); len(publish) == 0 {
			var err error
			if publish, err = i.exposedPorts(); err != nil {
				return nil, err
			}
		}
	}
	for i, pub := range publish {
		if strings.IndexByte(pub, ':') < 0 { // single port
			port := strings.TrimSuffix(pub, `/tcp`)
			port = strings.TrimSuffix(port, `/udp`)
			publish[i] = port + `:` + pub
		}
	}
	return publish, nil
}

func (i Image) exposedPorts() ([]string, error) {
	var m map[string]interface{}
	if output, err := cmd.Run(cmd.O{Output: true}, `docker`, `inspect`,
		`-f`, `{{ json .Config.ExposedPorts }}`, release.ImageNameOf(i.svcName),
	); err != nil {
		return nil, err
	} else if err := json.Unmarshal([]byte(output), &m); err != nil {
		return nil, err
	}
	ports := []string{}
	for k := range m {
		ports = append(ports, k)
	}
	return ports, nil
}
