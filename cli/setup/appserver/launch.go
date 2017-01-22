package appserver

import (
	"os"
	"os/exec"
	"path"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/fs"
	"github.com/bughou-go/xiaomei/utils/process"
	"github.com/fatih/color"
)

func Launch() {
	f := fs.OpenAppend(path.Join(config.App.Root(), `log/appserver.log`))
	defer f.Close()

	config.Logf(f, `starting.`)
	app, _ := cmd.Start(cmd.O{Panic: true, Stdout: f, Stderr: f}, `./`+config.App.Name())

	result := process.WaitPort(app.Process.Pid, config.App.Port(), config.App.StartTimeout(), false)
	switch result {
	case `ok`:
		config.Logf(f, color.GreenString(`started. (`+config.Servers.CurrentAppServer().AppAddr()+`)`))
		WaitApp(app, `crashed.`, f)
	case `died`:
		WaitApp(app, `starting failed.`, f)
	case `timeout`:
		app.Process.Kill()
		config.Logf(f, color.RedString(`starting timeout.`))
	default:
		app.Process.Kill()
		panic(`unknown result: ` + result)
	}
}

func WaitApp(app *exec.Cmd, msg string, f *os.File) {
	if err := app.Wait(); err != nil {
		config.Logf(f, color.RedString(msg+` (`+err.Error()+`)`))
	} else {
		config.Logf(f, color.RedString(msg))
	}
}
