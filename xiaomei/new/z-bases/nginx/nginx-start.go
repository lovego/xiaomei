package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/fatih/color"
)

func init() {
	go handleSignals()
}

func main() {
	confData := getConfData()
	if n := generateConf(confData); n > 0 {
		if names := getNamesFromAddrs(confData.BackendAddrs); len(names) > 0 {
			waitUntilNamesResolved(names)
		}
	}
	log(color.GreenString(`started. (:%s)`, confData.ListenPort))
	cmd := exec.Command(`nginx`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

type configData struct {
	ListenPort   string
	BackendAddrs []string
}

func getConfData() configData {
	port := os.Getenv(`NGPORT`)
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
	}
}

func generateConf(confData configData) int {
	tmplFiles, err := filepath.Glob(`/etc/nginx/sites-available/*.conf.tmpl`)
	if err != nil {
		panic(err)
	}
	for _, tmplFile := range tmplFiles {
		confFile := strings.Replace(tmplFile, `available`, `enabled`, 1)
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

func handleSignals() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	s := <-c
	println(` killed by ` + s.String() + ` signal.`)
	os.Exit(0)
}

const ISO8601 = `2006-01-02T15:04:05Z0700`

func log(msg interface{}) {
	fmt.Println(time.Now().Format(ISO8601), msg)
}
