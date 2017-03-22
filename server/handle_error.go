package server

import (
	"time"

	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/server/xm"
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
	if res.Size() <= 0 {
		res.Json(map[string]string{"code": "404", "message": "Not Found."})
	}
}

func handleServerError(req *xm.Request, res *xm.Response) {
	res.WriteHeader(500)
	if res.Size() <= 0 {
		res.Json(map[string]string{"code": "500", "message": "Application Server Error."})
	}
}
