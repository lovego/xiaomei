package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/fatih/color"
)

func main() {
	confData := getConfData()
	if n := generateConf(confData); n > 0 {
		if names := getNamesFromAddrs(confData.BackendAddrs); len(names) > 0 {
			waitUntilNamesResolved(names)
		}
	}
	log(color.GreenString(`started. (:%s)`, confData.ListenPort))
	if err := syscall.Exec(`/usr/sbin/nginx`, []string{`nginx`}, nil); err != nil {
		panic(err)
	}
}

type configData struct {
	ListenPort   string
	BackendAddrs []string
	SendfileOff  bool
}

func getConfData() configData {
	port := os.Getenv(`NGINXPORT`)
	if port == `` {
		port = `8000`
	}
	var addrs []string
	if addrsStr := os.Getenv(`NGBackendAddrs`); addrsStr != `` {
		addrs = strings.Split(addrsStr, `,`)
	}
	return configData{
		ListenPort:   port,
		BackendAddrs: addrs,
		SendfileOff:  os.Getenv(`SendfileOff`) == `true`,
	}
}

func generateConf(confData configData) int {
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
	return len(tmplFiles)
}

func makeConf(file string, confData configData) []byte {
	var buf bytes.Buffer
	if err := template.Must(template.ParseFiles(file)).Execute(&buf, confData); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func getNamesFromAddrs(addrs []string) (names []string) {
	for _, addr := range addrs {
		if host, _, err := net.SplitHostPort(addr); err != nil {
			panic(err)
		} else if net.ParseIP(host) == nil {
			names = append(names, host)
		}
	}
	return
}

func waitUntilNamesResolved(names []string) {
	fmt.Printf("%d names: %s.\n", len(names), strings.Join(names, `, `))

	start := time.Now()
	var wg sync.WaitGroup
	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			waitNameResolved(name)
			wg.Done()
		}(name)
	}
	wg.Wait()

	duration := time.Since(start)
	duration -= duration % time.Second
	color.Green("%d names: all ready in %s.\n\n", len(names), duration)
}

var debug = os.Getenv(`debug`) != ``

func waitNameResolved(host string) {
	start := time.Now()
	sleep := time.Second
	for {
		if _, err := net.LookupHost(host); err == nil {
			break
		} else {
			if debug {
				fmt.Println(err, sleep)
			}
			time.Sleep(sleep)
			sleep += sleep / 10
		}
	}
	duration := time.Since(start)
	duration -= duration % time.Second
	fmt.Printf("%s ready in %s.\n", host, duration)
}

const ISO8601 = `2006-01-02T15:04:05Z0700`

func log(msg interface{}) {
	fmt.Println(time.Now().Format(ISO8601), msg)
}
