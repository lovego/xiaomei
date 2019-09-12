package middlewares

import (
	"net/url"
	"strings"

	"github.com/lovego/config"
	"github.com/lovego/goa/middlewares"
)

var CORS = middlewares.NewCORS(allowOrigin)

func allowOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	hostname := u.Hostname()
	return strings.HasSuffix(hostname, config.Domain()) || hostname == `localhost`
}
