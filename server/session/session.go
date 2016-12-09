package session

import (
	"net/http"
)

type Store interface {
	Get(req *http.Request, p interface{})
	Set(res http.ResponseWriter, data interface{})
}
