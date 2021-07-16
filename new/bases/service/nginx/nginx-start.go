package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/fatih/color"
)

func main() {
	port := os.Getenv("ProPORT")
	if port == "" {
		port = "8000"
	}
	addr := ":" + port

	process := startNginx(port)
	log.Printf("starting.(%s)", addr)
	waitPortReady(process.Pid, addr)
	log.Println(color.GreenString("frontend started. (%s)", addr))

	// SIGUSR1 for log reopen
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		if err := process.Signal(<-c); err != nil {
			log.Println(err)
		}
	}
}

func startNginx(port string) *os.Process {
	generateConf(port)

	cmd := exec.Command("nginx")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		} else {
			log.Println("shutdown")
			os.Exit(0)
		}
	}()

	return cmd.Process
}

func generateConf(port string) {
	config := struct {
		ListenPort  string
		SendfileOff bool
	}{
		ListenPort:  port,
		SendfileOff: os.Getenv("SendfileOff") == "true",
	}
	tmplFiles, err := filepath.Glob("/etc/nginx/sites-available/*.conf.tmpl")
	if err != nil {
		log.Panic(err)
	}
	for _, tmplFile := range tmplFiles {
		var buf bytes.Buffer
		if err := template.Must(template.ParseFiles(tmplFile)).Execute(&buf, config); err != nil {
			log.Panic(err)
		}
		confFile := `/etc/nginx/sites-enabled/` + strings.TrimSuffix(filepath.Base(tmplFile), `.tmpl`)
		if err := ioutil.WriteFile(confFile, buf.Bytes(), 0644); err != nil {
			log.Panic(err)
		}
	}
}

func waitPortReady(pid int, addr string) {
	for i := 0; i < 7; i++ {
		cmd := exec.Command("lsof", "-aP", "-i"+addr, "-stcp:listen")
		cmd.Stderr = os.Stderr
		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		if len(out) > 0 {
			fmt.Print(string(out))
			return
		}
		time.Sleep(time.Second)
	}
	log.Printf("waitPortReady timeout(%s)\n", addr)
	os.Exit(1)
}
