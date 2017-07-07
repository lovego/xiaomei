package main

import (
	"testing"
)

func TestMakeConf(t *testing.T) {
	println(string(makeConf(`z_test.tmpl`, configData{
		ListenPort:   `8001`,
		BackendAddrs: []string{`127.0.0.1:3001`, `127.0.0.1:3002`},
	})))
}
