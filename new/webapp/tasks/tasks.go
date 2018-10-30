package tasks

import (
	"context"

	"github.com/lovego/config"
	"github.com/lovego/config/db/redisdb"
	loggerPkg "github.com/lovego/logger"
	"github.com/lovego/redis-cron"
)

var logger = config.NewLogger("cron.log")

func Start() {
	c := cron.New(redisdb.Pool("default"))
	if err := c.AddFunc("0 * * * * *", exampleTask); err != nil {
		logger.Panic(err)
	}
	c.Start()
}

func exampleTask() {
	logger.Record(func(ctx context.Context) error {
		// work to do goes here
		return nil
	}, nil, func(l *loggerPkg.Fields) {
		l.With("type", "example")
	})
}
