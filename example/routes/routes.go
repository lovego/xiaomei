package routes

import (
	"time"

	"github.com/bughou-go/xiaomei/server"
)

type session struct {
	UserId    int
	UserName  string
	LoginTime time.Time
}

func Get() *server.Router {
	router := server.NewRouter()

	router.Get(`/`, func(req *server.Request, res *server.Response) {
		res.Json(map[string]string{`hello`: `world`})
	})

	router.Get(`/session-get`, func(req *server.Request, res *server.Response) {
		var sess = session{}
		req.Session(&sess)
		res.Json(sess)
	})

	router.Get(`/session-set`, func(req *server.Request, res *server.Response) {
		var sess = session{UserId: 1, UserName: `xiaomei`, LoginTime: time.Now()}
		res.Session(sess)
		res.Json(sess)
	})

	return router
}
