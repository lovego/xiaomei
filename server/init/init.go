// this package must be imported first,
// to ensure it is initialized before all the other package.
package init

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/utils"
)

func init() {
	utils.Log(fmt.Sprintf(`starting.(%s)`, config.Env()))
	go handleSignals()
}

func handleSignals() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	s := <-c
	println(` killed by ` + s.String() + ` signal.`)
	os.Exit(0)
}
