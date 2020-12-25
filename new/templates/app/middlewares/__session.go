package middlewares

import (
	"github.com/lovego/goa"
	"github.com/lovego/tracer"
	"{{ .ProPath }}/middlewares/helpers"
)

func SessionParse(c *goa.Context) {
	parseSession(c)
	c.Next()
}

func parseSession(c *goa.Context) {
	ck, _ := c.Request.Cookie(helpers.Cookie.Name)
	if ck == nil || ck.Value == "" {
		return
	}
	ck.MaxAge = helpers.Cookie.MaxAge

	var data helpers.Session
	if err := helpers.CookieStore.Get(ck, &data); err != nil {
		tracer.Log(c.Context(), "session: ", err)
		return
	}

	c.Set("session", data)
}
