package app

import (
	"errors"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/images/app/db"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return append([]*cobra.Command{
		{
			Use:   `build-bin`,
			Short: `build the app server.`,
			RunE:  release.NoArgCall(buildBinary),
		},
		DepsCmd(),
		copy2vendorCmd(),
	}, db.Cmds()...)
}

func buildBinary() error {
	utils.Log(color.GreenString(`building app binary.`))
	if cmd.Ok(cmd.O{
		Dir: filepath.Join(release.Root(), `..`),
		Env: []string{`GOBIN=` + release.App().Root()},
	}, `go`, `install`, `-v`) {
		return nil
	}
	return errors.New(`building app binary failed.`)
}
