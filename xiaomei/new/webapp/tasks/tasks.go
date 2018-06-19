package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/fatih/color"
	loggerPkg "github.com/lovego/logger"
	"github.com/lovego/xiaomei/config"
	"github.com/robfig/cron"
)

var logger = config.NewLogger("cron.log")
var debug = os.Getenv("debugTasks") != ""

func main() {
	runtime.GOMAXPROCS(1)

	c := cron.New()
	if err := c.AddFunc("0 * * * * *", exampleTask); err != nil {
		panic(err)
	}
	c.Start()
	log.Println(color.GreenString(`started.`))
	select {}
}

func exampleTask() {
	logger.Record(debug, func(ctx context.Context) error {
		// work to do goes here
		return nil
	}, nil, func(l *loggerPkg.Fields) {
		l.With("type", "example")
	})
}
