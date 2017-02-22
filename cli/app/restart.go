package appserver

import (
	"bytes"
	"errors"
	"path"
	"strings"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/process"
)

func Restart(daemon bool) error {
	// stop current
	Stop()
	// remove current
	Remove()
	// start new
	return Start(daemon)
}

func Running() string {
	buf := bytes.Buffer{}
	output, _ := cmd.Run(cmd.O{Output: true, Stderr: &buf},
		`docker`, `inspect`, `--type=container`, `--format={{ .State.Running }}`, config.Deploy.Name(),
	)

	if stdErr := buf.String(); stdErr != `` &&
		stdErr != "Error: No such container: "+config.Deploy.Name()+"\n" {
		print(stdErr)
	}
	return strings.TrimSpace(output)
}

func Stop() {
	if Running() == `true` {
		out, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `docker`, `stop`, config.Deploy.Name())
		if out != config.Deploy.Name() {
			println(out)
		}
	}
}

func Remove() {
	if Running() != `` {
		out, _ := cmd.Run(cmd.O{Panic: true, Output: true}, `docker`, `rm`, config.Deploy.Name())
		if out != config.Deploy.Name() {
			println(out)
		}
	}
}

func Start(daemon bool) error {
	tail := cmd.TailFollow(path.Join(config.App.Root(), `log/appserver.log`))
	defer tail.Process.Kill()

	StartDocker(daemon)
	if daemon && process.WaitPort(getAppServerPid(),
		config.App.Port(), config.App.StartTimeout()+3*time.Second, true,
	) != `ok` {
		cmd.Run(cmd.O{Panic: true}, `docker`, `stop`, config.Deploy.Name())
		return errors.New(`start AppServer failed.`)
	}
	// wait one more step, ensure tail have gotten the started message.
	time.Sleep(process.WaitStep)
	return nil
}

func StartDocker(daemon bool) {
	root := config.App.Root()
	xiaomei, _ := cmd.Run(cmd.O{Output: true, Panic: true}, `which`, `xiaomei`)
	args := []string{
		`run`, `--name=` + config.Deploy.Name(),
		`-v`, root + `:` + root,
		`-v`, xiaomei + `:/usr/local/bin/xiaomei`,
		`-w`, root, `--network=host`,
	}
	if daemon {
		args = append(args, `-d`, `--restart=always`)
	} else {
		args = append(args, `--rm`)
	}
	args = append(args, config.App.DockerImage(), `xiaomei`, `launch`)

	cmd.Run(cmd.O{Panic: true}, `docker`, args...)
}

func getAppServerPid() int {
	deployName := config.Deploy.Name()
	ppid, _ := cmd.Run(cmd.O{Output: true, Panic: true},
		`docker`, `inspect`, `--format`, `{{ .State.Pid }}`, deployName)
	if ppid = strings.TrimSpace(ppid); ppid == `` {
		panic(`empty AppServer ppid.`)
	}
	pid := process.ChildPid(ppid, config.App.Name(), 3*time.Second)
	if pid <= 0 {
		panic(`find appserver pid failed.`)
	}
	return pid
}
