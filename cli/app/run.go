package app

import (
	"errors"
	"path/filepath"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Run() error {
	if err := buildBinary(); err != nil {
		return err
	}
	tail := cmd.TailFollow(
		filepath.Join(config.App.Root(), `log/app.log`),
		filepath.Join(config.App.Root(), `log/app.err`),
	)
	defer tail.Process.Kill()

	startDocker()
	return nil
}

func startDocker() {
	rootDir := config.App.Root()
	cmd.Run(cmd.O{Panic: true}, `docker`,
		`run`, `--name=`+config.Cluster.DeployName(), `-it`, `--rm`, `--network=host`,
		`-v`, rootDir+`:/home/ubuntu/appserver`,
		`-v`, rootDir+`/log:/home/ubuntu/appserver/log`,
		config.App.DockerImage(),
	)
}

func buildBinary() error {
	config.Log(`building.`)
	if cmd.Ok(cmd.O{Env: []string{`GOBIN=` + config.App.Root()}}, `go`, `install`) {
		return nil
	}
	return errors.New(`build failed.`)
}
