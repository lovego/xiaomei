package session

import (
	"net/http"
)

type Session interface {
	Get(req *http.Request, p interface{})
	Set(req *http.Request, res http.ResponseWriter, data interface{})
}

type CookieSession struct {
	Name string
}
