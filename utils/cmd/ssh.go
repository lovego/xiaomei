package cmd

import (
	"errors"
	"net"
	// "os/exec"
	"os/user"
	"strings"

	"github.com/bughou-go/xiaomei/utils/slice"
)

func SshRun(o O, addr, shellScript string) (string, error) {
	isLocal, err := IsLocalAddr(addr)
	if err != nil {
		return ``, err
	}
	if isLocal {
		return Run(o, `bash`, `-c`, shellScript)
	} else {
		return Run(o, `ssh`, `-t`, addr, shellScript)
	}
}

func IsLocalAddr(addr string) (bool, error) {
	i := strings.IndexByte(addr, '@')
	if i <= 0 {
		return false, errors.New(`invalid addr: ` + addr)
	}
	if ok, err := IsCurrentUser(addr[0:i]); !ok || err != nil {
		return ok, err
	}
	return IsLocalHost(addr[i+1:])
}

func IsCurrentUser(name string) (bool, error) {
	u, err := user.Current()
	if err != nil {
		return false, err
	}
	return u != nil && (u.Username == name || u.Name == name), nil
}

func IsLocalHost(addr string) (bool, error) {
	ips, err := net.LookupIP(addr)
	if err != nil {
		return false, err
	}
	for _, ip := range ips {
		if ip.IsLoopback() {
			return true, nil
		}
	}
	machineAddrs := MachineAddrs()
	for _, ip := range ips {
		if slice.ContainsString(machineAddrs, ip.String()) {
			return true, nil
		}
	}
	return false, nil
}

var theMachineAddrs []string

func MachineAddrs() []string {
	if theMachineAddrs != nil {
		return theMachineAddrs
	}

	ifcAddrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	addrs := make([]string, len(ifcAddrs))
	for i, ifcAddr := range ifcAddrs {
		addr := ifcAddr.String()
		if i := strings.IndexByte(addr, '/'); i >= 0 {
			addr = addr[:i]
		}
		addrs[i] = addr
	}
	theMachineAddrs = addrs
	return theMachineAddrs
}
