package godoc

import (
	"errors"
	"os"

	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/utils/cmd"
)

const (
	defaultDomain        = `godoc.dev`
	defaultAddr          = `0.0.0.0:1234`
	defaultIndexInterval = `1s`
)

func InDevelop() error {
	inProject := config.InProject()
	if inProject && config.Godoc.Port() == `` {
		return nil
	}
	if err := setupUpstartInDevelop(inProject); err != nil {
		return err
	}

	return setupNginxInDevelop(inProject)
}

func setupUpstartInDevelop(inProject bool) error {
	var gopath, godoc, addr, indexInterval string
	if gopath = os.Getenv(`GOPATH`); gopath == `` {
		return errors.New(`empty GOPATH.`)
	}
	godoc, _ = cmd.Run(cmd.O{Output: true, Panic: true}, `which`, `godoc`)
	if inProject {
		addr = config.Servers.CurrentAppServer().GodocAddr()
		indexInterval = config.Godoc.IndexInterval()
	} else {
		addr = defaultAddr
		indexInterval = defaultIndexInterval
	}

	return setupUpstart(&upstartConf{
		GoPath:        gopath,
		GodocBin:      godoc,
		Addr:          addr,
		IndexInterval: indexInterval,
	})
}

func setupNginxInDevelop(inProject bool) error {
	var domain string
	var addrs []string

	if inProject {
		domain = config.Godoc.Domain()
		for _, server := range config.Servers.All() {
			if server.HasTask(`appserver`) {
				addrs = append(addrs, server.GodocAddr())
			}
		}
	} else {
		domain = defaultDomain
		addrs = []string{defaultAddr}
	}

	return setupNginx(&nginxConf{
		Addrs:  addrs,
		Domain: domain,
	})
}
