package process

import (
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lovego/xiaomei/utils/cmd"
)

const WaitStep = 100 * time.Millisecond

// wait until the process has bound the port.
func WaitPort(pid int, port string, timeout time.Duration, sudo bool) string {
	binary := `lsof`
	args := []string{`-ap`, strconv.Itoa(pid), `-itcp:` + port}
	if sudo {
		args = append([]string{binary}, args...)
		binary = `sudo`
	}

	for w := time.Duration(0); w <= timeout; w += WaitStep {
		if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true}, binary, args...) {
			return `ok`
		}
		if !Alive(pid) {
			return `died`
		}
		time.Sleep(WaitStep)
	}
	return `timeout`
}

func Alive(pid int) bool {
	return syscall.Kill(pid, syscall.Signal(0)) == nil
}

var prefixDigits = regexp.MustCompile(`^\d+`)
var suffixDigits = regexp.MustCompile(`^\d+`)

func ChildPid(ppid, child string, timeout time.Duration) int {
	for w := time.Duration(0); w <= timeout; w += WaitStep {
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
		time.Sleep(WaitStep)
	}
	return -1
}
