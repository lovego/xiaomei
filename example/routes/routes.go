package routes

import (
	"time"

	"github.com/bughou-go/xiaomei/server/xm"
)

type session struct {
	UserId    int
	UserName  string
	LoginTime time.Time
}

func Get() *xm.Router {
	router := xm.NewRouter()

	router.Get(`/`, func(req *xm.Request, res *xm.Response) {
		res.Json(map[string]string{`hello`: `world`})
	})

	router.Get(`/session-get`, func(req *xm.Request, res *xm.Response) {
		var sess = session{}
		req.Session(&sess)
		res.Json(sess)
	})

	router.Get(`/session-set`, func(req *xm.Request, res *xm.Response) {
		var sess = session{UserId: 1, UserName: `xiaomei`, LoginTime: time.Now()}
		res.Session(sess)
		res.Json(sess)
	})

	return router
}
