package deploy

import (
	"github.com/lovego/xiaomei/xiaomei/deploy/host"
	// "github.com/lovego/xiaomei/xiaomei/deploy/swarm"
)

type driver interface {
	Deploy(svcName string) error
	Logs(svcName string, all bool) error
	Ps(svcName string, watch bool, options []string) error
	// register to release package
	ImageNameOf(svcName string) string
	PortsOf(svcName string) []string
	ServiceNames() map[string]bool
	// called by Run
	FlagsForRun(svcName string) ([]string, error)
}

var theDriver driver

func getDriver() driver {
	if theDriver == nil {
		theDriver = host.Driver
	}
	return theDriver
}

func Deploy(svcName string) error {
	return getDriver().Deploy(svcName)
}

func Logs(svcName string, all bool) error {
	return getDriver().Logs(svcName, all)
}

func Ps(svcName string, watch bool, options []string) error {
	return getDriver().Ps(svcName, watch, options)
}

func ImageNameOf(svcName string) string {
	return getDriver().ImageNameOf(svcName)
}

func PortsOf(svcName string) []string {
	return getDriver().PortsOf(svcName)
}

func ServiceNames() map[string]bool {
	return getDriver().ServiceNames()
}
