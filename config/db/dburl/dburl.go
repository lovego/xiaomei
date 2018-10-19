package dburl

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

type URL struct {
	URL              *url.URL
	MaxIdle, MaxOpen int
	MaxLife          time.Duration
}

func Parse(str string) URL {
	uri, err := url.Parse(str)
	if err != nil {
		log.Panic(err)
	}
	u := URL{URL: uri}

	q := u.URL.Query()
	if str := q.Get("maxIdle"); str != "" {
		q.Del("maxIdle")
		u.MaxIdle = parseInt(str)
	} else if os.Getenv("GOENV") == "production" {
		u.MaxIdle = 1
	} else {
		u.MaxIdle = 0
	}

	if str := q.Get("maxOpen"); str != "" {
		q.Del("maxOpen")
		u.MaxOpen = parseInt(str)
	} else {
		u.MaxOpen = 10
	}

	if str := q.Get("maxLife"); str != "" {
		q.Del("maxLife")
		u.MaxLife = parseDuration(str)
	} else {
		u.MaxLife = 10 * time.Minute
	}
	u.URL.RawQuery = q.Encode()
	return u
}

func parseInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Panic(err)
	}
	return value
}

func parseDuration(str string) time.Duration {
	value, err := time.ParseDuration(str)
	if err != nil {
		log.Panic(err)
	}
	return value
}
