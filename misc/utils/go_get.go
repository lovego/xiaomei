package utils

import (
	"github.com/lovego/cmd"
)

func GoGetMod(args ...string) error {
	_, err := cmd.Run(
		cmd.O{
			Env: []string{`GOPROXY=https://goproxy.cn,direct`, `GO111MODULE=on`},
		},
		`go`, append([]string{`get`, `-v`}, args...)...,
	)
	return err
}
