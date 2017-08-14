package deploy

import (
	"github.com/lovego/xiaomei/xiaomei/deploy/simple"
	// "github.com/lovego/xiaomei/xiaomei/deploy/swarm"
)

type driver interface {
	FlagsForRun(svcName string) ([]string, error)
	Deploy(svcName, feature string) error
	RmDeploy(svcName, feature string) error
	Logs(svcName, feature string) error
	Ps(svcName, feature string, watch bool, options []string) error
}

var theDriver driver

func getDriver() driver {
	if theDriver == nil {
		theDriver = simple.Driver
	}
	return theDriver
}
