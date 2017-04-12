package deploy

import (
	"github.com/lovego/xiaomei/xiaomei/deploy/simple"
	// "github.com/lovego/xiaomei/xiaomei/deploy/swarm"
)

type driver interface {
	FlagsForRun(svcName string) ([]string, error)
	AccessAddrs(svcName string) []string
	Deploy(svcName string) error
	RmDeploy(svcName string) error
	Logs(svcName string) error
	Ps(svcName string, watch bool, options []string) error
}

var theDriver driver

func getDriver() driver {
	if theDriver == nil {
		theDriver = simple.Driver
	}
	return theDriver
}
