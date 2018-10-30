package filter

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/config"
)

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
	if !isAllowed(origin) {
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
		res.Header().Set(`Access-Control-Allow-Headers`,
			`X-Requested-With, Content-Type, withCredentials`)
		return true
	}
	return false
}

func isAllowed(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	hostname := u.Hostname()
	if strings.HasSuffix(hostname, config.Domain()) || hostname == `localhost` {
		return true
	}
	return false
}
