package routes

import (
	"errors"
	"time"

	"github.com/lovego/config"
	"github.com/lovego/errs"
	"github.com/lovego/goa"
	"github.com/lovego/tracer"
	"{{ .ProPath }}/helpers"
)

func Setup(router *goa.Router)  {
	router.Get(`/`, func(c *goa.Context) {
		tracer.Tag(c.Context(), `hello`, `world`)
		c.Json(map[string]string{`hello`: config.DeployName()})
	})

	router.Get(`/data`, func(c *goa.Context) {
		c.Data([]string{`hello`, `world`}, nil)
	})

	router.Get(`/business-error`, func(c *goa.Context) {
		c.Data([]string{`hello`, `world`}, errs.New(`captcha-err`, `验证码错误`))
	})

	router.Get(`/other-error`, func(c *goa.Context) {
		c.Data([]string{`hello`, `world`}, errs.Trace(errors.New(`unknown error`)))
	})

	router.Get(`/session-get`, func(c *goa.Context) {
		c.Json(helpers.GetSession(c))
	})

	router.Get(`/session-set`, func(c *goa.Context) {
		sess := helpers.Session{UserId: 100, UserName: `xiaomei`, LoginTime: time.Now()}
		helpers.SaveSession(c, sess)
		c.Write([]byte(`ok`))
	})

	router.Get(`/session-delete`, func(c *goa.Context) {
		helpers.DeleteSession(c)
		c.Write([]byte(`ok`))
	})
}
