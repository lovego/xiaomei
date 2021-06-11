package release

import (
	"os"

	"github.com/lovego/cmd"
)

func GoCmd() string {
	if cmd := os.Getenv("GoCMD"); cmd != "" {
		return cmd
	}
	return `go`
}

func GoGetByProxy(args ...string) error {
	_, err := cmd.Run(
		cmd.O{
			Env: []string{`GOPROXY=https://goproxy.cn,direct`, `GO111MODULE=on`},
		},
		GoCmd(), append([]string{`get`, `-v`}, args...)...,
	)
	return err
}
