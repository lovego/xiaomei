package godoc

import (
	"errors"
	"path/filepath"

	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/utils/install"
)

func InDeploy() error {
	if config.Godoc.Port() == `` {
		return nil
	}
	if err := SetupUpstartInDeploy(); err != nil {
		return err
	}
	return SetupNginxInDeploy()
}

func SetupUpstartInDeploy() error {
	if config.Godoc.Port() == `` {
		return nil
	}
	deployRoot := config.Deploy.Root()
	if _, err := cmd.Run(cmd.O{}, `sudo`, `ln`, `-Tsf`, `.`, filepath.Join(deployRoot, `src`)); err != nil {
		return err
	}
	godocBin, _ := cmd.Run(cmd.O{Output: true}, `which`, `godoc`)
	if godocBin == `` {
		if err := install.AptGet(`golang-go.tools`); err != nil {
			return err
		}
		godocBin, _ = cmd.Run(cmd.O{Output: true}, `which`, `godoc`)
	}
	if godocBin == `` {
		return errors.New(`godoc not found.`)
	}

	return setupUpstart(&upstartConf{
		GoPath:        deployRoot,
		GodocBin:      godocBin,
		Addr:          config.Servers.CurrentAppServer().GodocAddr(),
		IndexInterval: config.Godoc.IndexInterval(),
	})
}

func SetupNginxInDeploy() error {
	if config.Godoc.Port() == `` {
		return nil
	}
	var addrs []string
	for _, server := range config.Servers.All() {
		if server.HasTask(`appserver`) {
			addrs = append(addrs, server.GodocAddr())
		}
	}
	return setupNginx(&nginxConf{
		Domain: config.Godoc.Domain(),
		Addrs:  addrs,
	})
}
