package appserver

import (
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
)

func Wait() {
	ppidStr := getPpid()
	if ppidStr == `` {
		return
	}
	pidStr := getPid(ppidStr)
	if pidStr == `` {
		return
	}

	WaitPort(parseInt(ppidStr), parseInt(pidStr))
}

// wait until the AppPort has been bound.
func WaitPort(ppid, pid int) {
	pidStr := strconv.Itoa(pid)
	const step = 100 * time.Millisecond
	for w := time.Duration(0); w <= config.App.StartTimeout(); w += step {
		if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true},
			`lsof`, `-ap`, pidStr, `-itcp:`+config.App.Port(),
		) {
			config.Log(color.GreenString(`started. (` + config.Servers.CurrentAppServer().AppAddr() + `)`))
			return
		}
		if !processAlive(ppid) || !processAlive(pid) {
			config.Log(color.RedString(`starting failed.`))
			return
		}
		time.Sleep(step)
	}
	syscall.Kill(-ppid, syscall.SIGTERM) // kill process group
	config.Log(`starting timeout.`)
}

func processAlive(pid int) bool {
	return syscall.Kill(pid, syscall.Signal(0)) == nil
}

var prefixDigits = regexp.MustCompile(`^\d+`)
var suffixDigits = regexp.MustCompile(`^\d+`)

func getPpid() string {
	deployName := config.Deploy.Name()
	status, _ := cmd.Run(cmd.O{Output: true, Panic: true}, `status`, `apps/`+deployName)
	prefix := `apps/` + deployName + ` start/post-start, process `
	if !strings.HasPrefix(status, prefix) {
		config.Log(`unexpected status: ` + status + `.`)
		return ``
	}
	ppidStr := status[len(prefix):]
	ppidStr = prefixDigits.FindString(ppidStr)
	if ppidStr == `` {
		config.Log(`unexpected status: ` + status + `.`)
		return ``
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
	config.Log(`find appserver pid timeout.`)
	return ``
}

func parseInt(str string) int {
	if i, err := strconv.Atoi(str); err != nil {
		panic(err)
	} else {
		return i
	}
}
