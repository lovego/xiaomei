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
	var wg sync.WaitGroup
	for _, addr := range backends {
		wg.Add(1)
		go func(addr string) {
			waitTcpReady(addr)
			wg.Done()
		}(addr)
	}
	wg.Wait()
	fmt.Printf("all %d backends ready(%ds).\n", len(backends), time.Since(start)/time.Second)
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
	fmt.Printf("%s ready(%ds).\n", addr, time.Since(start)/time.Second)
}

func getBackends() (result []string) {
	if paths, err := filepath.Glob(`/etc/nginx/sites-enabled/*.conf`); err != nil {
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
