package server

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/server/xm"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/fs"
)

var accessLog, errLog = setupLogger()

func setupLogger() (*os.File, *os.File) {
	if config.DevMode() {
		return os.Stdout, os.Stderr
	} else {
		return fs.OpenAppend(filepath.Join(config.Root(), `log/app.log`)),
			fs.OpenAppend(filepath.Join(config.Root(), `log/app.err`))
	}
}

func writeLog(req *xm.Request, res *xm.Response, t time.Time, err interface{}) []byte {
	line := getLogLine(req, res, t, err)
	if err != nil {
		errLog.Write(line)
	} else {
		accessLog.Write(line)
	}
	return line
}

func getLogLine(req *xm.Request, res *xm.Response, t time.Time, err interface{}) []byte {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ' '
	writer.Write(getLogFields(req, res, t, err))
	writer.Flush()
	return buf.Bytes()
}

/*
  $time_iso8601 $host $request_method $request_uri $content_length $server_protocol
  $status $body_bytes_sent
  $request_time
  $session $remote_addr $http_referer $http_user_agent, $error, $stack
*/
func getLogFields(req *xm.Request, res *xm.Response, t time.Time, err interface{}) []string {
	slice := []string{t.Format(utils.ISO8601), req.Host,
		req.Method, req.URL.RequestURI(), strconv.FormatInt(req.ContentLength, 10), req.Proto,
		strconv.FormatInt(res.Status(), 10), strconv.FormatInt(res.Size(), 10),
		fmt.Sprintf(`%.6f`, time.Since(t).Seconds()),
		fmt.Sprint(req.Session), req.ClientAddr(), req.Referer(), req.UserAgent(),
	}
	if err != nil {
		slice = append(slice, fmt.Sprint(err), string(utils.Stack(6)))
	}
	for i, v := range slice {
		v = strings.TrimSpace(v)
		if v == `` {
			v = `-`
		}
		slice[i] = v
	}
	return slice
}
