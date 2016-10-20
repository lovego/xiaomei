package config

import (
	"net"
	"os"
	"path"
	"strings"
	"time"
)

var Root = rootDir()
var Data = parseConfigData()
var TimeZone = time.FixedZone(Data.TimeZoneName, Data.TimeZoneOffset)

func init() {
	setupMailer()
}

func rootDir() string {
	program := os.Args[0]
	if !path.IsAbs(program) {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		program = path.Join(cwd, program)
	}
	if strings.HasPrefix(program, `/tmp/`) || strings.HasSuffix(program, `.test`) {
		// only for development
		gopath := os.Getenv(`GOPATH`)
		if gopath != `` {
			return path.Join(gopath, `src/github.com/bughou-go/xiaomei/release/`)
		} else {
			panic(`detect root dir failed.`)
		}
	} else {
		feature := `config/config.yml`
		dir := path.Dir(program)
		for ; dir != `/`; dir = path.Dir(dir) {
			if _, err := os.Stat(path.Join(dir, feature)); err == nil {
				return dir
			}
		}
		panic(`detect root dir failed.`)
	}
}

func CurrentAppServer() ServerConfig {
	ifcAddrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, server := range Data.DeployServers {
		if server.AppAddr != `` {
			for _, ifcAddr := range ifcAddrs {
				if strings.HasPrefix(ifcAddr.String(), server.Addr+`/`) {
					return server
				}
			}
		}
	}
	return ServerConfig{}
}
