package cmd

import (
	"fmt"
	"strings"
)

const SshShareConnFlags = `-o ControlMaster=auto` +
	` -o ControlPath=/tmp/ssh_mux_%h_%p_%r` +
	` -o ControlPersist=600`

func SshRun(o O, addr, script string, flags ...string) (output string, err error) {
	args := strings.Split(SshShareConnFlags, ` `)
	args = append(args, flags...)
	args = append(args, addr, script)
	return Run(o, `ssh`, args...)
}

func SshJumpRun(o O, jumpAddr, addr, script string) (output string, err error) {
	if script == `` {
		return SshRun(o, jumpAddr, fmt.Sprintf(
			`ssh -t %s %s`, SshShareConnFlags, addr,
		), `-t`)
	}
	if o.PrintCmd() {
		fmt.Println(script)
	}
	if _, err := SshRun(O{Stdin: strings.NewReader(script)}, jumpAddr, fmt.Sprintf(
		`ssh %s %s 'cat > /tmp/runScript.sh'`, SshShareConnFlags, addr,
	)); err != nil {
		return ``, err
	}
	return SshRun(o, jumpAddr, fmt.Sprintf(
		`ssh -t %s %s bash /tmp/runScript.sh`, SshShareConnFlags, addr,
	), `-t`)
}
