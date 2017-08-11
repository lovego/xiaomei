package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/lovego/xiaomei/utils/httputil"
)

type Container struct {
	Id    string
	Names []string
}

func main() {
	signal, names := getArgs()
	client := getClient()
	containers := getContainers(names, client)
	for _, c := range containers {
		killContainer(c, signal, client)
	}
}

func killContainer(c Container, signal string, client *httputil.Client) {
	queryUrl := `/containers/` + c.Id + `/kill?` + url.Values{`signal`: {signal}}.Encode()
	resp, err := client.Post(queryUrl, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := resp.Check(http.StatusNoContent); err != nil {
		log.Fatal(err)
	}
	log.Println(c.Names)
}

func getContainers(names []string, client *httputil.Client) []Container {
	filters := map[string][]string{`name`: names, `status`: {`running`}}
	buf, err := json.Marshal(filters)
	if err != nil {
		log.Fatal(err)
	}
	queryUrl := `/containers/json?` + url.Values{`filters`: {string(buf)}}.Encode()

	var data []Container
	if err := client.GetJson(queryUrl, nil, nil, &data); err != nil {
		log.Fatal(err)
	}
	return data
}

func getClient() *httputil.Client {
	return &httputil.Client{
		BaseUrl: `http://docker/v1.30`,
		Client: &http.Client{
			Transport: &http.Transport{
				Dial: unixDial,
			},
			Timeout: 5 * time.Second,
		},
	}
}

func unixDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", "/var/run/docker.sock")
}

func getArgs() (signal string, names []string) {
	flag.CommandLine.Usage = usage
	flag.StringVar(&signal, `s`, `SIGKILL`, `signal to send to the containers`)
	flag.Parse()

	names = flag.Args()
	if len(names) == 0 {
		usage()
		os.Exit(0)
	}
	return
}

func usage() {
	fmt.Fprintf(os.Stderr,
		`docker-kill send signal to containers whose name match the regexps. (version 17.8.11)
usage: %s [options] <container-name-regexp> ...
`, os.Args[0])
	flag.PrintDefaults()
}
