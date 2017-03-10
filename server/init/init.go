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
	go logSignals()
}

func logSignals() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	s := <-c
	println(`stopped. (` + s.String() + `)`)
	os.Exit(0)
}
