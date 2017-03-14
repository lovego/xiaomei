// this package must be imported first,
// to ensure it is initialized before all the other package.
package init

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bughou-go/xiaomei/utils"
)

func init() {
	utils.Log(`starting.`)
	go handleSignals()
}

func handleSignals() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	s := <-c
	println(` killed by ` + s.String() + ` signal.`)
	os.Exit(0)
}
