package main

import (
	_ "github.com/lovego/config/init" // this package must be the first.

	"runtime"

	"github.com/lovego/goa"
	middles "github.com/lovego/goa/middlewares"
	"github.com/lovego/goa/server"
	"github.com/lovego/goa/utilroutes"
	"{{ .ProPath }}/middlewares"
	"{{ .ProPath }}/routes"
	// "{{ .ProPath }}/tasks"
)

func main() {
	if n := runtime.NumCPU() - 1; n >= 1 {
		runtime.GOMAXPROCS(n)
	}

	router := goa.New()
	router.Use(middlewares.Logger.Record)
	middles.SetupProcessingList(router)
	router.Use(middlewares.SessionParse)
	router.Use(middlewares.CORS.Check)

	utilroutes.Setup(router)
	routes.Setup(router)

	// tasks.Start()
	server.ListenAndServe(router)
}
