package deploy

import (
	"github.com/lovego/xiaomei/xiaomei/deploy/stack"
)

type driver interface {
	Deploy(svcName string, rmCurrent bool) error
	Logs(svcName string, all bool) error
	Ps(svcName string, watch bool, options []string) error
}

var Driver driver = stack.Driver
