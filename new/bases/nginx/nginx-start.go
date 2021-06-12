package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/fatih/color"
)

func main() {
	port := os.Getenv(`ProPORT`)
	if port == `` {
		port = `8000`
	}
	if conn, _ := net.DialTimeout("tcp", port, time.Second); conn != nil {
		conn.Close()
		log.Printf("port %s is adready bound by other process.", port)
		return
	}

	startNginx(port)
	waitPortReady(":" + port)

	log.Println(color.GreenString(`started. (:%s)`, port))
	select {}
}

func startNginx(port string) *exec.Cmd {
	generateConf(port)

	cmd := exec.Command(`/usr/sbin/nginx`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		if err := cmd.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()

	return cmd
}

type configData struct {
	ListenPort  string
	SendfileOff bool
}

func generateConf(port string) {
	tmplFiles, err := filepath.Glob(`/etc/nginx/sites-available/*.conf.tmpl`)
	if err != nil {
		panic(err)
	}
	confData := configData{
		ListenPort:  port,
		SendfileOff: os.Getenv(`SendfileOff`) == `true`,
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

func waitPortReady(port string) {
	for {
		time.Sleep(100 * time.Millisecond)
		if conn, _ := net.DialTimeout("tcp", port, time.Second); conn != nil {
			conn.Close()
			return
		}
	}
}
