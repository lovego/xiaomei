package main

import (
	_ "github.com/lovego/config/init" // this package must be the first.

	"os"
	"path/filepath"
	"runtime"
	"context"

	"github.com/lovego/config"
	"github.com/lovego/config/db/redisdb"
	"github.com/lovego/logger"
	"github.com/lovego/fs"
	"github.com/lovego/goa"
	"github.com/lovego/goa/server"
	"github.com/lovego/goa/utilroutes"
	"github.com/lovego/redis-cron"
	"{{ .ModulePath }}/core"
	"{{ .ModulePath }}/generic"
	"{{ .ModulePath }}/generic/middlewares"
	"{{ .ModulePath }}/support"
)

func main() {
	if n := runtime.NumCPU() - 1; n >= 1 {
		runtime.GOMAXPROCS(n)
	}

	router := goa.New()
	router.Use(middlewares.CORS.Check, middlewares.Logger.Record, middlewares.SessionParse)
	if os.Getenv("GOA_DOC") != `` {
		router.DocDir(filepath.Join(fs.SourceDir(), "..", "docs", "apis"))
		setupRoutes(&router.RouterGroup)
		os.Exit(0)
	}
	setupRoutes(&router.RouterGroup)

	// setupCrontab()
	server.ListenAndServe(router)
}

func setupRoutes(router *goa.RouterGroup) {
	utilroutes.Setup(router)
	router.Use(middlewares.Filter)

	router.Get(`/`, func(c *goa.Context) {
		c.Json(map[string]string{`hello`: config.DeployName()})
	})

	core.Routes(router)
	generic.Routes(router)
	support.Routes(router)
}

var cronLogger = config.NewLogger("cron.log")

func setupCrontab() {
	crontab := cron.New(redisdb.Pool("default"))
	if err := crontab.AddFunc("0 * * * * *", exampleTask); err != nil {
		cronLogger.Panic(err)
	}
	crontab.Start()
}

func exampleTask() {
	cronLogger.Record(func(ctx context.Context) error {
		// work to do goes here
		return nil
	}, nil, func(l *logger.Fields) {
		l.With("type", "example")
	})
}
