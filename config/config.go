package config

import (
	"net"
	"strings"
	"time"
)

var Root = rootDir()
var Data = parseConfigData()
var TimeZone = time.FixedZone(Data.TimeZoneName, Data.TimeZoneOffset)

func init() {
	setupMailer()
}

func CurrentAppServer() ServerConfig {
	ifc_addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, server := range Data.DeployServers {
		if server.AppAddr != `` {
			for _, ifc_addr := range ifc_addrs {
				if strings.HasPrefix(ifc_addr.String(), server.Addr+`/`) {
					return server
				}
			}
		}
	}
	return ServerConfig{}
}
