package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lovego/config"
	"github.com/lovego/config/conf"
	"github.com/lovego/errs"
	"github.com/lovego/goa"
	"github.com/lovego/tracer"
	"github.com/lovego/xiaomei/new/webapp/helpers"
)

func Filter(c *goa.Context) {
	if err := check(c); err != nil {
		c.Data(nil, err)
		return
	}

	c.Next()
}

func check(c *goa.Context) error {
	path := c.Request.URL.Path
	if isPublicPaths(path) {
		return nil
	}
	if strings.HasPrefix(path, "/sign/") {
		err := checkSign(c.Request.Header)
		if err != nil {
			tracer.Tag(c.Context(), `timestamp`, c.Request.Header.Get(`Timestamp`))
			tracer.Tag(c.Context(), `sign`, c.Request.Header.Get(`Sign`))
		}
		return err
	}
	return checkSession(c)
}

func isPublicPaths(path string) bool {
	return path == "/"
}

func checkSign(header http.Header) error {
	timestamp, sign := header.Get(`Timestamp`), header.Get(`Sign`)
	if err := checkTimestamp(timestamp); err != nil {
		return err
	}
	if conf.TimestampSign(timestamp, config.Secret()) != sign {
		return errs.New("sign-err", "header Sign error")
	}
	return nil
}

func checkTimestamp(timestamp string) error {
	stamp, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return errs.New("args-err", "header Timestamp error")
	}
	diff := time.Now().Unix() - stamp
	if diff > 60 || diff < -60 {
		return errs.New("args-err", "header Timestamp has a gap more than one minute")
	}
	return nil
}

func checkSession(c *goa.Context) error {
	if helpers.GetSession(c).UserId <= 0 {
		return errs.New(`token-err`, `token error, please login again`)
	}
	return nil
}
