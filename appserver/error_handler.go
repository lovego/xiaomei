package main

import (
	"fmt"
	"log"
	//	"github.com/bughou-go/xiaomei/filter"
	"github.com/bughou-go/xiaomei/config"
	"strconv"
	"strings"
	"time"
	"github.com/bughou-go/xm"
)

func errorHandler(t time.Time, req *xm.Request, res *xm.Response) {
	err := recover()
	if err != nil {
		serverError(req, res)
	}
	access := accessLog(req, res, t)
	if err != nil {
		msg := fmt.Sprintf("%s\nPANIC: %s\n%s", access, err, xm.Stack(4))
		go config.AlarmMail(strconv.FormatInt(res.Status(), 10)+"错误", msg)
		log.Print(msg)
	}
	fmt.Println(access)
}

func accessLog(req *xm.Request, res *xm.Response, t time.Time) string {
	return strings.Join([]string{
		t.Format(`2006-01-02 15:04:05 -0700`),
		req.Host, req.Method, req.URL.RequestURI(), req.ClientAddr(), fmt.Sprint(req.Session),
		strconv.FormatInt(res.Status(), 10), strconv.FormatInt(res.Size(), 10),
		time.Since(t).String(), req.Referer(), req.UserAgent(),
	}, ` `)
}

func notFound(req *xm.Request, res *xm.Response) {
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

func serverError(req *xm.Request, res *xm.Response) {
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
