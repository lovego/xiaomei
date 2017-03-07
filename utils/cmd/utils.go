package cmd

import (
	"os/exec"
)

func TailFollow(paths ...string) *exec.Cmd {
	Run(O{Panic: true}, `touch`, append([]string{`-a`}, paths...)...)
	tail, err := Start(O{Panic: true}, `tail`, append([]string{`-fqn0`}, paths...)...)
	if err != nil {
		panic(err)
	}
	return tail
}
