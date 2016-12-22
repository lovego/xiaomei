package appserver

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Wait() {
	ppidStr := getPpid()
	pidStr := getPid(ppidStr)
	ppid := parseInt(ppidStr)
	pid := parseInt(pidStr)

	// wait until the AppPort has been bound.
	for w := time.Duration(0); w <= config.App.StartTimeout(); w += time.Second {
		if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true},
			`lsof`, `-ap`, pidStr, `-itcp:`+config.App.Port(),
		) {
			exit(`started. (` + config.Servers.Current().AppAddr + `:` + config.App.Port() + `)`)
		}
		if !processAlive(ppid) || !processAlive(pid) {
			exit(`starting failed.`)
		}
		time.Sleep(time.Second)
	}
	syscall.Kill(-ppid, syscall.SIGTERM) // kill process group
	exit(`starting timeout.`)
}

func processAlive(pid int) bool {
	return syscall.Kill(pid, syscall.Signal(0)) == nil
}

var prefixDigits = regexp.MustCompile(`^\d+`)
var suffixDigits = regexp.MustCompile(`^\d+`)

func getPpid() string {
	deployName := config.Deploy.Name()
	status, _ := cmd.Run(cmd.O{Output: true, Panic: true}, `status`, deployName)
	prefix := deployName + ` start/post-start, process `
	if !strings.HasPrefix(status, prefix) {
		exit(`unexpected status: ` + status + `.`)
	}
	ppidStr := status[len(prefix):]
	ppidStr = prefixDigits.FindString(ppidStr)
	if ppidStr == `` {
		exit(`unexpected status: ` + status + `.`)
	}
	return ppidStr
}

func getPid(ppid string) string {
	for i := 0; i < 10; i++ {
		output, _ := cmd.Run(cmd.O{Output: true}, `ps`, `-o`, `pid,cmd`, `--no-headers`, `--ppid`, ppid)
		index := strings.Index(output, ` ./`+config.App.Name())
		if index > 0 {
			pidStr := suffixDigits.FindString(strings.TrimSpace(output[0:index]))
			if pidStr != `` {
				return pidStr
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	exit(`find appserver pid timeout.`)
	return ``
}

func parseInt(str string) int {
	if i, err := strconv.Atoi(str); err != nil {
		panic(err)
	} else {
		return i
	}
}

func exit(msg string) {
	println(time.Now().Format(config.ISO8601), msg)
	os.Exit(0)
}
