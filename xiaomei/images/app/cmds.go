package app

import (
	"errors"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lovego/utils"
	"github.com/lovego/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/images/app/db"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return append([]*cobra.Command{
		{
			Use:   `compile`,
			Short: `compile the app server binary.`,
			RunE:  release.NoArgCall(compile),
		},
		specCmd(),
		depsCmd(),
		copy2vendorCmd(),
	}, db.Cmds()...)
}

func compile() error {
	utils.Log(color.GreenString(`compile the app server binary.`))
	if cmd.Ok(cmd.O{
		Dir: filepath.Join(release.Root(), `..`),
		Env: []string{`GOBIN=` + filepath.Join(release.Root(), `img-app`)},
	}, `go`, `install`, `-v`) {
		return gospec(nil, true)
	}
	return errors.New(`compile the app server binary failed.`)
}
