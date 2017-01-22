package process

import (
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bughou-go/xiaomei/utils/cmd"
)

// wait until the process has bound the port.
func WaitPort(pid int, port string, timeout time.Duration) string {
	pidStr := strconv.Itoa(pid)
	const step = 100 * time.Millisecond
	for w := time.Duration(0); w <= timeout; w += step {
		if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true}, `lsof`, `-ap`, pidStr, `-itcp:`+port) {
			return `ok`
		}
		if !Alive(pid) {
			return `died`
		}
		time.Sleep(step)
	}
	return `timeout`
}

func Alive(pid int) bool {
	return syscall.Kill(pid, syscall.Signal(0)) == nil
}

var prefixDigits = regexp.MustCompile(`^\d+`)
var suffixDigits = regexp.MustCompile(`^\d+`)

func ChildPid(ppid, child string, timeout time.Duration) int {
	const step = 100 * time.Millisecond
	for w := time.Duration(0); w <= timeout; w += step {
		output, _ := cmd.Run(cmd.O{Output: true}, `ps`, `-o`, `pid,cmd`, `--no-headers`, `--ppid`, ppid)
		if index := strings.Index(output, ` ./`+child); index > 0 {
			if str := suffixDigits.FindString(strings.TrimSpace(output[0:index])); str != `` {
				if pid, err := strconv.Atoi(str); err != nil {
					panic(err)
				} else {
					return pid
				}
			}
		}
		time.Sleep(step)
	}
	return -1
}
