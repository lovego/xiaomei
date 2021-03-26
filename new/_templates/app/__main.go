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
	// "{{ .ModulePath }}/tasks"
)

func main() {
	if n := runtime.NumCPU() - 1; n >= 1 {
		runtime.GOMAXPROCS(n)
	}

	router := goa.New()
	router.Use(middlewares.CORS.Check)
	router.Use(middlewares.Logger.Record)
	router.Use(middlewares.SessionParse)
	utilroutes.Setup(router)
	router.Use(middlewares.Filter)

	if os.Getenv("GOA_DOC") != `` {
		router.DocDir(filepath.Join(fs.SourceDir(), "docs", "apis"))
		routes.Setup(router)
		return
	}
	routes.Setup(router)

	// tasks.Start()
	server.ListenAndServe(router)
}
