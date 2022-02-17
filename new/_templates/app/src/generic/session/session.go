package session

import (
	"time"

	"github.com/lovego/config"
	"github.com/lovego/goa"
	"github.com/lovego/sessions/cookiestore"
	"github.com/lovego/tracer"
)

var Cookie = config.Cookie()
var CookieStore = cookiestore.New(config.Secret())

type Session struct {
	UserId    int
	UserName  string
	LoginTime time.Time
}

func Get(c *goa.Context) Session {
	v := c.Get("session")
	if data, ok := v.(Session); ok {
		return data
	}
	return Session{}
}

func Save(c *goa.Context, data Session) {
	err := CookieStore.Save(c.ResponseWriter, &Cookie, data)
	if err != nil {
		tracer.Log(c.Context(), "save session: ", err)
	}
}

func Delete(c *goa.Context) {
	CookieStore.Delete(c.ResponseWriter, &Cookie)
}
