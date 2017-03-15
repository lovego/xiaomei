package server

import (
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/server/xm"
)

func handleError(t time.Time, req *xm.Request, res *xm.Response, notFound *bool) {
	if *notFound {
		handleNotFound(req, res)
	}

	err := recover()
	if err != nil {
		handleServerError(req, res)
	}
	if err == nil && req.URL.Path == alivePath {
		return
	}
	log := writeLog(req, res, t, err)
	if err != nil {
		go config.Alarm(`500错误`, string(log))
	}
}

func handleNotFound(req *xm.Request, res *xm.Response) {
	res.WriteHeader(404)
	if res.Size() > 0 {
		return
	}
	if req.Header.Get("X-Requested-With") != "" {
		res.Write([]byte(`{ "message": "404" }`))
	} else {
		res.Render(`error/404`, nil)
	}
}

func handleServerError(req *xm.Request, res *xm.Response) {
	res.WriteHeader(500)
	if res.Size() > 0 {
		return
	}
	if req.Header.Get("X-Requested-With") != "" {
		res.Write([]byte(`{ "message": "500" }`))
	} else {
		res.Render(`error/500`, nil)
	}
}
