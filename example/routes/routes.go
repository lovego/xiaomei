package routes

import (
	"github.com/bughou-go/xiaomei/server"
)

func Get() *server.Router {
	router := server.NewRouter()

	router.Get(`/`, func(req *server.Request, res *server.Response) {
		res.Json(map[string]string{`hello`: `world`})
	})

	return router
}
