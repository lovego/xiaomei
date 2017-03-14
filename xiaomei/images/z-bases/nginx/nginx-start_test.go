package main

import (
	"testing"
)

func TestWaitTcpReady(t *testing.T) {
	waitTcpReady(`none:http`)
	waitTcpReady(`localhost:50012`)
}
