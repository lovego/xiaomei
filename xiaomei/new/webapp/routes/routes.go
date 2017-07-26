package routes

import (
	"time"

	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/router"
	"github.com/lovego/xiaomei/utils/errs"
)

type session struct {
	UserId    int
	UserName  string
	LoginTime time.Time
}

func Routes() *router.Router {
	router := router.New()

	router.Get(`/`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Json(map[string]string{`hello`: config.DeployName()})
		req.Log(map[string]interface{}{`hello`: `world`, `i'm`: `xiaomei`})
	})

	router.Get(`/result-ok`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Result([]string{`hello`, `world`}, nil)
	})

	router.Get(`/result-bad`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Result([]string{`hello`, `world`}, errs.New(`captcha-err`, `验证码错误`))
	})

	router.Get(`/session-get`, func(req *xiaomei.Request, res *xiaomei.Response) {
		var sess = session{}
		req.Session(&sess)
		res.Json(sess)
	})

	router.Get(`/session-set`, func(req *xiaomei.Request, res *xiaomei.Response) {
		var sess = session{UserId: 1, UserName: `xiaomei`, LoginTime: time.Now()}
		res.Session(sess)
		res.Json(sess)
	})

	router.Get(`/session-delete`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Session(nil)
	})

	return router
}
