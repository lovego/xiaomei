package appserver

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Restart() {
	// stop current
	Stop()
	// start new
	Start()
}

func Stop() {
	var buf bytes.Buffer
	cmd.Run(cmd.O{Stderr: &buf}, `sudo`, `stop`, `apps/`+config.Deploy.Name())
	stdErr := buf.String()
	if stdErr != "stop: Unknown instance: \n" {
		print(stdErr)
	}
}

func Start() {
	tail := tailLog()
	defer tail.Process.Kill()
	output, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `sudo`, `start`, `apps/`+config.Deploy.Name())

	println(output)
	if !strings.Contains(output, `start/running,`) {
		os.Exit(1)
	}
}

func tailLog() *exec.Cmd {
	appserverLog := path.Join(config.App.Root(), `log/appserver.log`)
	cmd.Run(cmd.O{Panic: true}, `touch`, `-a`, appserverLog)
	tail, err := cmd.Start(cmd.O{Panic: true}, `tail`, `-fn0`, appserverLog)
	if err != nil {
		panic(err)
	}
	return tail
}
