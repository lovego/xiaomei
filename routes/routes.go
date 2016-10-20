package routes

import (
	"github.com/bughou-go/xm"
)

func Get() *xm.Router {
	var router = xm.NewRouter()

	router.Get(`/`, func(req *xm.Request, res *xm.Response) {
		res.Json(map[string]string{`hello`: `world`})
	})

	return router
}
