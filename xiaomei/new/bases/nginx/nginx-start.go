package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/fatih/color"
)

func main() {
	confData := getConfData()
	generateConf(confData)
	log(color.GreenString(`started. (:%s)`, confData.ListenPort))
	if err := syscall.Exec(`/usr/sbin/nginx`, []string{`nginx`}, nil); err != nil {
		panic(err)
	}
}

type configData struct {
	ListenPort  string
	SendfileOff bool
}

func getConfData() configData {
	port := os.Getenv(`NGINXPORT`)
	if port == `` {
		port = `8000`
	}
	return configData{
		ListenPort:  port,
		SendfileOff: os.Getenv(`SendfileOff`) == `true`,
	}
}

func generateConf(confData configData) {
	tmplFiles, err := filepath.Glob(`/etc/nginx/sites-available/*.conf.tmpl`)
	if err != nil {
		panic(err)
	}
	for _, tmplFile := range tmplFiles {
		confFile := `/etc/nginx/sites-enabled/` + strings.TrimSuffix(filepath.Base(tmplFile), `.tmpl`)
		if err := ioutil.WriteFile(confFile, makeConf(tmplFile, confData), 0644); err != nil {
			panic(err)
		}
	}
}

func makeConf(file string, confData configData) []byte {
	var buf bytes.Buffer
	if err := template.Must(template.ParseFiles(file)).Execute(&buf, confData); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

const ISO8601 = `2006-01-02T15:04:05Z0700`

func log(msg interface{}) {
	fmt.Println(time.Now().Format(ISO8601), msg)
}
