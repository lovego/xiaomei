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
	confData := getConfData()
	generateConf(confData)

	cmd := exec.Command(`/usr/sbin/nginx`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	go func() {
		if err := cmd.Run(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()

	waitPort(":" + confData.ListenPort)
	log.Println(color.GreenString(`started. (:%s)`, confData.ListenPort))

	select {}
}

func waitPort(port string) {
	for {
		if conn, _ := net.DialTimeout("tcp", port, time.Second); conn != nil {
			conn.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
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
