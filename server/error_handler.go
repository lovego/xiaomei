package server

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
)

func handleError(t time.Time, req *Request, res *Response, notFound *bool) {
	if *notFound {
		handleNotFound(req, res)
	}
	err := recover()
	if err != nil {
		handleServerError(req, res)
	}
	access := accessLog(req, res, t)
	if err != nil {
		msg := fmt.Sprintf("%s\nPANIC: %s\n%s", access, err, utils.Stack(4))
		go config.AlarmMail(strconv.FormatInt(res.Status(), 10)+"错误", msg)
		log.Print(msg)
	}
	fmt.Println(access)
}

func accessLog(req *Request, res *Response, t time.Time) string {
	var session interface{}
	req.Session(&session)

	return strings.Join([]string{
		t.Format(`2006-01-02 15:04:05 -0700`),
		req.Host, req.Method, req.URL.RequestURI(), req.ClientAddr(), fmt.Sprint(session),
		strconv.FormatInt(res.Status(), 10), strconv.FormatInt(res.Size(), 10),
		time.Since(t).String(), req.Referer(), req.UserAgent(),
	}, ` `)
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
