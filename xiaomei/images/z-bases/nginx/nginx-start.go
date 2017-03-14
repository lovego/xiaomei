package main

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

func main() {
	waitBackends()
	cmd := exec.Command(`nginx`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func waitBackends() {
	start := time.Now()
	backends := getBackends()
	fmt.Printf("%d backends: %s.\n", len(backends), strings.Join(backends, `, `))
	var wg sync.WaitGroup
	for _, addr := range backends {
		wg.Add(1)
		go func(addr string) {
			waitTcpReady(addr)
			wg.Done()
		}(addr)
	}
	wg.Wait()
	duration := time.Since(start)
	duration -= duration % time.Second
	color.New(color.FgGreen).Printf("%d backends: all ready in %s.\n\n", len(backends), duration)
}

var debug = os.Getenv(`debug`) != ``

func waitTcpReady(addr string) {
	sleep := time.Second
	start := time.Now()
	for {
		if _, err := net.Dial(`tcp`, addr); err == nil {
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
	fmt.Printf("%s ready in %s.\n", addr, duration)
}

func getBackends() (result []string) {
	if paths, err := filepath.Glob(`/etc/nginx/sites-enabled/*`); err != nil {
		panic(err)
	} else {
		for _, p := range paths {
			if addrs := backendAddrs(p); len(addrs) > 0 {
				result = append(result, addrs...)
			}
		}
	}
	return result
}

func backendAddrs(p string) (result []string) {
	f, err := os.Open(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if addr := parseBackendAddr(scanner.Text()); addr != `` {
			result = append(result, addr)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return
}

var reProxyPass = regexp.MustCompile(`^\s*proxy_pass\s+(http://\S+)\s*;`)

func parseBackendAddr(line string) (addr string) {
	matches := reProxyPass.FindStringSubmatch(line)
	if len(matches) < 2 {
		return ``
	}

	if uri, err := url.Parse(matches[1]); err != nil {
		panic(err)
	} else {
		addr = uri.Host
	}
	if strings.IndexByte(addr, ':') < 0 {
		addr += `:http`
	}
	return
}
