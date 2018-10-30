package middlewares

import (
	"github.com/lovego/goa"
	"github.com/lovego/tracer"
	"github.com/lovego/xiaomei/new/webapp/helpers"
)

func ParseSession(c *goa.Context) {
	parseSession(c)
	c.Next()
}

func parseSession(c *goa.Context) {
	ck, err := c.Request.Cookie(helpers.Cookie.Name)
	if err != nil {
		tracer.Log(c.Context(), "get session cookie error:", err)
		return
	}
	if ck == nil || ck.Value == "" {
		return
	}
	ck.MaxAge = helpers.Cookie.MaxAge

	var data helpers.Session
	if err := helpers.CookieStore.Get(ck, &data); err != nil {
		tracer.Log(c.Context(), "get session error:", err)
		return
	}

	c.Set("session", data)
}
