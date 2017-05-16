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

	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/server/xm"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/fs"
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

func writeLog(
	req *xm.Request, res *xm.Response, t time.Time, err bool, errStr, stack string,
) []byte {
	line := getLogLine(req, res, t, err, errStr, stack)
	if err {
		errLog.Write(line)
	} else {
		accessLog.Write(line)
	}
	return line
}

func getLogLine(
	req *xm.Request, res *xm.Response, t time.Time, err bool, errStr, stack string,
) []byte {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ' '
	writer.Write(getLogFields(req, res, t, err, errStr, stack))
	writer.Flush()
	return buf.Bytes()
}

/*
  $time_iso8601 $host $request_method $request_uri $content_length $server_protocol
  $status $body_bytes_sent
  $request_time
  $session $remote_addr $http_referer $http_user_agent, $error, $stack
*/
func getLogFields(
	req *xm.Request, res *xm.Response, t time.Time, err bool, errStr, stack string,
) []string {
	slice := []string{t.Format(utils.ISO8601), req.Host,
		req.Method, req.URL.RequestURI(), strconv.FormatInt(req.ContentLength, 10), req.Proto,
		strconv.FormatInt(res.Status(), 10), strconv.FormatInt(res.Size(), 10),
		fmt.Sprintf(`%.6f`, time.Since(t).Seconds()),
		getSession(req), req.ClientAddr(), req.Referer(), req.UserAgent(),
	}
	if err {
		slice = append(slice, errStr, stack)
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

func getSession(req *xm.Request) string {
	var sess interface{}
	req.Session(&sess)
	str := fmt.Sprint(sess)
	if strings.HasPrefix(str, `map[`) && strings.HasSuffix(str, `]`) {
		str = strings.TrimPrefix(str, `map[`)
		str = strings.TrimSuffix(str, `]`)
	}
	return str
}
