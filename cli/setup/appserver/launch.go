package appserver

import (
	"os/exec"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/process"
	"github.com/fatih/color"
)

func Launch() {
	color.NoColor = false

	config.Log(`starting.`)
	app, _ := cmd.Start(cmd.O{Panic: true}, `./`+config.App.Name())

	result := process.WaitPort(app.Process.Pid, config.App.Port(), config.App.StartTimeout(), false)
	switch result {
	case `ok`:
		config.Log(color.GreenString(`started. (` + config.Servers.CurrentAppServer().AppAddr() + `)`))
		WaitApp(app, `crashed.`)
	case `died`:
		WaitApp(app, `starting failed.`)
	case `timeout`:
		app.Process.Kill()
		config.Log(color.RedString(`starting timeout.`))
	default:
		app.Process.Kill()
		panic(`unknown result: ` + result)
	}
}

func WaitApp(app *exec.Cmd, msg string) {
	if err := app.Wait(); err != nil {
		config.Log(color.RedString(msg + ` (` + err.Error() + `)`))
	} else {
		config.Log(color.RedString(msg))
	}
}
