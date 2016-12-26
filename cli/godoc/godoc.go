package godoc

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/bughou-go/xiaomei/cli/install"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Start() error {
	if config.Godoc.Port() == `` {
		return nil
	}
	gopath := os.Getenv(`GOPATH`)
	if gopath == `` {
		return errors.New(`empty GOPATH.`)
	}
	godoc, _ := cmd.Run(cmd.O{Output: true, Panic: true}, `which`, `godoc`)
	if err := Upstart(gopath, godoc); err != nil {
		return err
	}
	return Nginx()
}

func Setup() error {
	if config.Godoc.Port() == `` {
		return nil
	}
	gopath := config.Deploy.Root()
	if _, err := cmd.Run(cmd.O{}, `sudo`, `ln`, `-Tsf`, `.`, filepath.Join(gopath, `src`)); err != nil {
		return err
	}
	godoc, _ := cmd.Run(cmd.O{Output: true}, `which`, `godoc`)
	if godoc == `` {
		if err := install.AptGet(`golang-go.tools`); err != nil {
			return err
		}
		godoc, _ = cmd.Run(cmd.O{Output: true}, `which`, `godoc`)
	}
	if godoc == `` {
		return errors.New(`godoc not found.`)
	}

	if err := Upstart(gopath, godoc); err != nil {
		return err
	}
	return Nginx()
}
