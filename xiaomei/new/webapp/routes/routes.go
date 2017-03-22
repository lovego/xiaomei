package routes

import (
	"time"

	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/server/xm"
)

type session struct {
	UserId    int
	UserName  string
	LoginTime time.Time
}

func Routes() *xm.Router {
	router := xm.NewRouter()

	router.Get(`/`, func(req *xm.Request, res *xm.Response) {
		res.Json(map[string]string{`hello`: config.DeployName()})
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

	router.Get(`/session-delete`, func(req *xm.Request, res *xm.Response) {
		res.Session(nil)
	})

	return router
}
