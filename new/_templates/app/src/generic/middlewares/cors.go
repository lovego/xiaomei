package middlewares

import (
	"net/url"

	"github.com/lovego/goa"
	"github.com/lovego/goa/middlewares"
)

var CORS = middlewares.NewCORS(allowOrigin)

func allowOrigin(origin string, c *goa.Context) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	hostname := u.Hostname()
	return hostname == `localhost`
}
