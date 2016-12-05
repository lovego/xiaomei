package config

import (
	"net"
	"strings"
	"time"
)

var timeZone *time.Location

func TimeZone() *time.Location {
	if timeZone == nil {
		time.FixedZone(Data().TimeZoneName, Data().TimeZoneOffset)
	}
	return timeZone
}

func CurrentAppServer() ServerConfig {
	ifcAddrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, server := range Data().DeployServers {
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
