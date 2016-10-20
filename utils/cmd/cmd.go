package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// option
type O struct {
	Stdin                io.Reader
	Stdout, Stderr       io.Writer
	Print, Panic, Output bool
}

func Run(o O, name string, args ...string) (output string, err error) {
	if o.Print {
		fmt.Println(name, strings.Join(args, ` `))
	}

	cmd := exec.Command(name, args...)
	setupStdIO(cmd, o)

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

func setupStdIO(cmd *exec.Cmd, o O) {
	if o.Stdin != nil {
		cmd.Stdin = o.Stdin
	} else {
		cmd.Stdin = os.Stdin
	}
	if o.Stdout != nil {
		cmd.Stdout = o.Stdout
	} else if !o.Output {
		cmd.Stdout = os.Stdout
	}
	if o.Stderr != nil {
		cmd.Stderr = o.Stderr
	} else {
		cmd.Stderr = os.Stderr
	}
}

func SudoWriteFile(file string, reader io.Reader) {
	cmd := exec.Command(`sudo`, `cp`, `/dev/stdin`, file)
	cmd.Stdin = reader
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	Run(O{Panic: true}, `sudo`, `chmod`, `644`, file)
}
