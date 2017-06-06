package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// option
type O struct {
	Stdin                       io.Reader
	Stdout, Stderr              io.Writer
	NoStdin, NoStdout, NoStderr bool
	Dir                         string
	Env                         []string
	Attr                        *syscall.SysProcAttr
	Print                       bool
	Panic                       bool // only for Run And Start
	Output                      bool // only for Run
}

func (o O) PrintCmd() bool {
	return o.Print || os.Getenv(`PrintCmd`) == `true`
}

func Run(o O, name string, args ...string) (output string, err error) {
	cmd := makeCmd(o, name, args)

	if o.Output {
		var bytes []byte
		bytes, err = cmd.Output()
		output = strings.TrimSpace(string(bytes))
	} else {
		err = cmd.Run()
	}

	if o.Panic && err != nil {
		panic(err)
	}

	return
}

func Start(o O, name string, args ...string) (cmd *exec.Cmd, err error) {
	cmd = makeCmd(o, name, args)

	err = cmd.Start()
	if o.Panic && err != nil {
		panic(err)
	}

	return cmd, err
}

func State(o O, name string, args ...string) *os.ProcessState {
	cmd := makeCmd(o, name, args)
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			panic(err)
		}
	}
	return cmd.ProcessState
}

func Ok(o O, name string, args ...string) bool {
	s := State(o, name, args...)
	return s != nil && s.Success()
}

func Fail(o O, name string, args ...string) bool {
	s := State(o, name, args...)
	return s != nil && !s.Success()
}

func makeCmd(o O, name string, args []string) *exec.Cmd {
	if o.PrintCmd() {
		fmt.Println(name, strings.Join(args, ` `))
	}

	cmd := exec.Command(name, args...)
	if o.Dir != `` {
		cmd.Dir = o.Dir
	}
	if o.Env != nil {
		cmd.Env = append(os.Environ(), o.Env...)
	}
	if o.Attr != nil {
		cmd.SysProcAttr = o.Attr
	}
	setupStdIO(cmd, o)
	return cmd
}

func setupStdIO(cmd *exec.Cmd, o O) {
	if !o.NoStdin {
		if o.Stdin != nil {
			cmd.Stdin = o.Stdin
		} else {
			cmd.Stdin = os.Stdin
		}
	}
	if !o.NoStdout {
		if o.Stdout != nil {
			cmd.Stdout = o.Stdout
		} else if !o.Output {
			cmd.Stdout = os.Stdout
		}
	}
	if !o.NoStderr {
		if o.Stderr != nil {
			cmd.Stderr = o.Stderr
		} else {
			cmd.Stderr = os.Stderr
		}
	}
}
