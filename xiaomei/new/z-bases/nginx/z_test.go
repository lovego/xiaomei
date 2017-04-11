package main

import (
	"testing"
)

func TestMakeConf(t *testing.T) {
	println(string(makeConf(`z_test.tmpl`, configData{})))
}

func TestWaitNameResolved(t *testing.T) {
	waitNameResolved(`baidu.com`)
}
