package routes

import (
	"errors"
	"time"

	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/router"
	"github.com/lovego/utils/errs"
)

type session struct {
	UserId    int
	UserName  string
	LoginTime time.Time
}

func Routes() *router.Router {
	routes := router.New()

	routes.Get(`/`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Json(map[string]string{`hello`: config.DeployName()})
		req.Log(map[string]interface{}{`hello`: `world`, `i'm`: `xiaomei`})
	})

	routes.Get(`/result`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Result([]string{`hello`, `world`}, nil)
	})

	routes.Get(`/result2`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Result([]string{`hello`, `world`}, errs.New(`captcha-err`, `验证码错误`))
	})

	routes.Get(`/result-err`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Result([]string{`hello`, `world`}, errors.New(`unknown error`))
	})

	routes.Get(`/json-err`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Result([]string{`hello`, `world`}, errs.Trace(errors.New(`unknown error`)))
	})

	routes.Get(`/session-get`, func(req *xiaomei.Request, res *xiaomei.Response) {
		var sess = session{}
		req.Session(&sess)
		res.Json(sess)
	})

	routes.Get(`/session-set`, func(req *xiaomei.Request, res *xiaomei.Response) {
		var sess = session{UserId: 1, UserName: `xiaomei`, LoginTime: time.Now()}
		res.Session(sess)
		res.Json(sess)
	})

	routes.Get(`/session-delete`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Session(nil)
	})

	return routes
}
