package session

import (
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
)

type Session interface {
	Get(req *http.Request, p interface{})
	Set(res http.ResponseWriter, data interface{})
}

type CookieSession struct {
	cookie http.Cookie
	secure *securecookie.SecureCookie
}

func NewCookieSession(c http.Cookie, secret string) *CookieSession {
	c.Value = ``
	c.HttpOnly = true
	return &CookieSession{
		cookie: c,
		secure: securecookie.New([]byte(secret), nil),
	}
}

func (cs *CookieSession) Get(req *http.Request, p interface{}) {
	ck, _ := req.Cookie(cs.cookie.Name)
	if ck != nil && ck.Value != `` {
		if err := cs.secure.Decode(ck.Name, ck.Value, p); err != nil {
			log.Println(err)
		}
	}
}

func (cs *CookieSession) Set(res http.ResponseWriter, data interface{}) {
	ck := cs.cookie // make a copy
	if data == nil {
		ck.MaxAge = -1
	} else {
		if encoded, err := cs.secure.Encode(cs.cookie.Name, data); err == nil {
			ck.Value = encoded
		} else {
			log.Println(err)
			return
		}
	}
	http.SetCookie(res, &ck)
}
