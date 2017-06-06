package cmd

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func TailFollow(paths ...string) *exec.Cmd {
	Run(O{Panic: true}, `touch`, append([]string{`-a`}, paths...)...)
	tail, err := Start(O{Panic: true}, `tail`, append([]string{`-fqn0`}, paths...)...)
	if err != nil {
		panic(err)
	}
	return tail
}

func SudoWriteFile(file string, reader io.Reader) {
	dir := filepath.Dir(file)
	if !isDir(dir) {
		Run(O{Panic: true}, `sudo`, `mkdir`, `-p`, dir)
	}

	cmd := exec.Command(`sudo`, `cp`, `/dev/stdin`, file)
	cmd.Stdin = reader
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	Run(O{Panic: true}, `sudo`, `chmod`, `644`, file)
}

func isDir(p string) bool {
	fi, _ := os.Stat(p)
	return fi != nil && fi.IsDir()
}
