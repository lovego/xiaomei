package main

import (
	_ "github.com/lovego/config/init" // this package must be the first.

  "runtime"

	"github.com/lovego/goa"
	"github.com/lovego/goa/server"
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
  router.Use(middlewares.CORS.Check)
  router.Use(middlewares.ParseSession)
  routes.Setup(router)

  // tasks.Start()
	server.ListenAndServe(router)
}

