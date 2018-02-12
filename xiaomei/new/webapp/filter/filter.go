package filter

import (
	"net/http"

	"github.com/lovego/xiaomei"
)

var theAllowedOrigins = map[string]bool{
	`http://www.example.com`: true,
}

func Process(req *xiaomei.Request, res *xiaomei.Response) bool {
	if returnCORS(req, res) {
		return false
	}
	return true
}

func returnCORS(req *xiaomei.Request, res *xiaomei.Response) bool {
	origin := req.Header.Get(`Origin`)
	if origin == `` {
		return false
	}
	if !theAllowedOrigins[origin] {
		res.WriteHeader(http.StatusForbidden)
		res.Write([]byte(`origin not allowed.`))
		return true
	}

	res.Header().Set(`Access-Control-Allow-Origin`, origin)
	res.Header().Set(`Access-Control-Allow-Credentials`, `true`)
	res.Header().Set(`Vary`, `Accept-Encoding, Origin`)

	if req.Method == `OPTIONS` { // preflight 预检请求
		res.Header().Set(`Access-Control-Max-Age`, `86400`)
		res.Header().Set(`Access-Control-Allow-Methods`, `GET, POST, PUT, DELETE`)
		return true
	}
	return false
}
