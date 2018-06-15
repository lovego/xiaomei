package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lovego/fs"
	loggerPkg "github.com/lovego/logger"
	"github.com/lovego/slice"
	"github.com/lovego/xiaomei"
	"github.com/lovego/xiaomei/config"
)

var logger = loggerPkg.New(getLogWriter())
var logBodyMethods = []string{http.MethodPost, http.MethodDelete, http.MethodPut}

func getLogWriter() io.Writer {
	if config.DevMode() {
		return os.Stdout
	}
	file, err := fs.NewLogFile(filepath.Join(config.Root(), `log`, `app.log`))
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func logFields(f *loggerPkg.Fields, req *xiaomei.Request, res *xiaomei.Response, debug bool) {
	f.With("host", req.Host)
	f.With("method", req.Method)
	f.With("path", req.URL.Path)
	f.With("rawQuery", req.URL.RawQuery)
	f.With("query", req.URL.Query())
	f.With("status", res.Status())
	f.With("reqBodySize", req.ContentLength)
	f.With("resBodySize", res.Size())
	// 	f.With("proto", req.Proto)
	f.With("ip", req.ClientAddr())
	f.With("agent", req.UserAgent())
	f.With("refer", req.Referer())
	var sess interface{}
	req.Session(&sess)
	if sess != nil {
		f.With("session", sess)
	}

	if logBody && slice.ContainsString(logBodyMethods, req.Method) || debug {
		f.With("reqBody", string(req.GetBody()))
		f.With("resBody", string(res.GetBody()))
	}
}
