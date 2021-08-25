package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lovego/config"
	config2 "github.com/lovego/config/config"
	"github.com/lovego/errs"
	"github.com/lovego/goa"
	"github.com/lovego/tracer"
	"{{ .ModulePath }}/middlewares/helpers"
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
	ts, err := parseTimestamp(timestamp)
	if err != nil {
		return err
	}
	if config2.TimestampSign(ts, config.Secret()) != sign {
		return errs.New("sign-err", "header Sign error")
	}
	return nil
}

func parseTimestamp(timestamp string) (int64, error) {
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return 0, errs.New("args-err", "header Timestamp error")
	}
	diff := time.Now().Unix() - ts
	if diff > 60 || diff < -60 {
		return 0, errs.New("args-err", "header Timestamp has a gap more than one minute")
	}
	return ts, nil
}

func checkSession(c *goa.Context) error {
	if helpers.GetSession(c).UserId <= 0 {
		return errs.New(`token-err`, `token error, please login again`)
	}
	return nil
}
