// this package must be imported first,
// to ensure it is initialized before all the other package.
package init

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bughou-go/xiaomei/utils/fs"
)

func init() {
	log(`starting.`)
	go logSignals()
}

func log(msg string) {
	const iso8601 = `2006-01-02T15:04:05Z0700`
	println(time.Now().Format(iso8601), msg)
}

func logSignals() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	s := <-c
	log(`stopped. (` + s.String() + `)`)
	os.Exit(0)
}
