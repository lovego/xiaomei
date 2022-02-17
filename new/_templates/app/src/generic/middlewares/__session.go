package middlewares

import (
	"github.com/lovego/goa"
	"github.com/lovego/tracer"
	"{{ .ModulePath }}/generic/session"
)

func SessionParse(c *goa.Context) {
	parseSession(c)
	c.Next()
}

func parseSession(c *goa.Context) {
	ck, _ := c.Request.Cookie(session.Cookie.Name)
	if ck == nil || ck.Value == "" {
		return
	}
	ck.MaxAge = session.Cookie.MaxAge

	var data session.Session
	if err := session.CookieStore.Get(ck, &data); err != nil {
		tracer.Log(c.Context(), "session: ", err)
		return
	}

	c.Set("session", data)
}
