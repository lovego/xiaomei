package routes

import (
	"errors"
	"time"

	"github.com/lovego/errs"
	"github.com/lovego/tracer"
	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/router"
)

type session struct {
	UserId    int
	UserName  string
	LoginTime time.Time
}

func Routes() *router.Router {
	routes := router.New()

	routes.Get(`/`, func(req *xiaomei.Request, res *xiaomei.Response) {
		ctx := req.Context()
		tracer.Tag(ctx, `hello`, `world`)
		res.Json(map[string]string{`hello`: config.DeployName()})
	})

	routes.Get(`/data`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Data([]string{`hello`, `world`}, nil)
	})

	routes.Get(`/business-error`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Data([]string{`hello`, `world`}, errs.New(`captcha-err`, `验证码错误`))
	})

	routes.Get(`/other-error`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Data([]string{`hello`, `world`}, errs.Trace(errors.New(`unknown error`)))
	})

	routes.Get(`/json-err`, func(req *xiaomei.Request, res *xiaomei.Response) {
		res.Json2([]string{`hello`, `world`}, errs.Trace(errors.New(`unknown error`)))
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
