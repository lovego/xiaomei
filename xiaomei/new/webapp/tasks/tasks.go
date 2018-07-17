package tasks

import (
	"context"
	"os"

	loggerPkg "github.com/lovego/logger"
	"github.com/lovego/redis-cron"
	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/config/db/redisdb"
)

var logger = config.NewLogger("cron.log")
var debug = os.Getenv("debugTasks") != ""

func Run() {
	c := cron.New(redisdb.Pool("default"))
	if err := c.AddFunc("0 * * * * *", exampleTask); err != nil {
		logger.Panic(err)
	}
	c.Run()
}

func exampleTask() {
	logger.Record(debug, func(ctx context.Context) error {
		// work to do goes here
		return nil
	}, nil, func(l *loggerPkg.Fields) {
		l.With("type", "example")
	})
}
