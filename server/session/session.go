package session

import (
	"encoding/json"
	"net/http"
)

type Session interface {
	Get(req *http.Request, p interface{})
	Set(req *http.Request, res http.ResponseWriter, data interface{})
}

type CookieSession struct {
	http.Cookie
}

func (cs *CookieSession) Get(req *http.Request, p interface{}) {
	cookie, _ := req.Cookie(cs.Cookie.Name)
	if cookie != nil && cookie.Value != `` {
		json.Unmarshal()
	}
}

func (cs *CookieSession) Set(req *http.Request, res http.ResponseWriter, data interface{}) {
}

func (cs *CookieSession) verifySig() {
}

func (cs *CookieSession) genSig() {
}
