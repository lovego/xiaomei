package main

import (
	_ "github.com/lovego/config/init" // this package must be the first.

	"os"
	"path/filepath"
	"runtime"

	"github.com/lovego/fs"
	"github.com/lovego/goa"
	"github.com/lovego/goa/server"
	"github.com/lovego/goa/utilroutes"
	"{{ .ModulePath }}/middlewares"
	"{{ .ModulePath }}/routes"
	//"{{ .ModulePath }}/tasks"
)

func main() {
	if n := runtime.NumCPU() - 1; n >= 1 {
		runtime.GOMAXPROCS(n)
	}

	router := goa.New()
	router.Use(middlewares.CORS.Check, middlewares.Logger.Record, middlewares.SessionParse)
	setupRouter(&router.RouterGroup)

	// tasks.Start()
	server.ListenAndServe(router)
}

func setupRouter(routerGroup *goa.RouterGroup) {
	if os.Getenv("GOA_DOC") != `` {
		routerGroup.DocDir(filepath.Join(fs.SourceDir(), "docs", "apis"))
		routes.Setup(routerGroup)
		os.Exit(1)
	}

	utilroutes.Setup(routerGroup)
	routerGroup.Use(middlewares.Filter)
	routes.Setup(routerGroup)
}
