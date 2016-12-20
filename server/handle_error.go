package server

import (
	"time"

	"github.com/bughou-go/xiaomei/config"
)

func handleError(t time.Time, req *Request, res *Response, notFound *bool) {
	if *notFound {
		handleNotFound(req, res)
	}

	err := recover()
	if err != nil {
		handleServerError(req, res)
	}
	log := writeLog(req, res, t, err)
	if err != nil {
		go config.App.Alarm(`500错误`, string(log))
	}
}

func handleNotFound(req *Request, res *Response) {
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

func handleServerError(req *Request, res *Response) {
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
