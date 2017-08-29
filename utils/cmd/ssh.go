package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/lovego/xiaomei/utils/slice"
)

const SshShareConnFlags = `-o ControlMaster=auto` +
	` -o ControlPath=/tmp/ssh_mux_%h_%p_%r` +
	` -o ControlPersist=600`

func SshRun(o O, addr, script string, flags ...string) (output string, err error) {
	args := strings.Split(SshShareConnFlags, ` `)
	if o.Stdin == nil && !slice.ContainsString(flags, `-t`) {
		args = append(args, `-t`)
	}
	args = append(args, flags...)
	args = append(args, addr, script)
	return Run(o, `ssh`, args...)
}

func SshJumpRun(o O, jumpAddr, addr, script string) (output string, err error) {
	var ttyFlag string
	if o.Stdin == nil {
		ttyFlag = `-t`
	}

	if script == `` {
		return SshRun(o, jumpAddr, fmt.Sprintf(
			`ssh %s %s %s`, ttyFlag, SshShareConnFlags, addr,
		))
	}
	if o.PrintCmd() {
		fmt.Println(script)
	}
	tmpFile := fmt.Sprintf(`/tmp/xiaomei.%s.sh`, time.Now().Format(`2006-01-02T15:04:05.999999999`))
	if _, err := SshRun(O{Stdin: strings.NewReader(script)}, jumpAddr, fmt.Sprintf(
		`ssh %s %s 'cat > %s'`, SshShareConnFlags, addr, tmpFile,
	)); err != nil {
		return ``, err
	}
	return SshRun(o, jumpAddr, fmt.Sprintf(
		`ssh %s %s %s "bash %s; rm %s"`, ttyFlag, SshShareConnFlags, addr, tmpFile, tmpFile,
	))
}
