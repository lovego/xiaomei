package helpers

import (
	"net/http"
	"time"

	"github.com/lovego/config"
	"github.com/lovego/goa"
	"github.com/lovego/sessions/cookiestore"
	"github.com/lovego/tracer"
)

var Cookie = &http.Cookie{
	Name:   "session",
	MaxAge: 86400 * 30,
}
var CookieStore = cookiestore.New(config.Secret())

type Session struct {
	UserId    int
	UserName  string
	LoginTime time.Time
}

func GetSession(c *goa.Context) Session {
	v := c.Get("session")
	if data, ok := v.(Session); ok {
		return data
	}
	return Session{}
}

func SaveSession(c *goa.Context, data Session) {
	err := CookieStore.Save(c.ResponseWriter, Cookie, data)
	if err != nil {
		tracer.Log(c.Context(), "save session: ", err)
	}
}

func DeleteSession(c *goa.Context) {
	CookieStore.Delete(c.ResponseWriter, Cookie)
}
