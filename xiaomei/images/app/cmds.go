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

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `app`,
		Short: `the app server.`,
	}
	cmd.AddCommand(&cobra.Command{
		Use:   `build-bin`,
		Short: `build the app server.`,
		RunE:  release.NoArgCall(buildBinary),
	})
	cmd.AddCommand(DepsCmd(), copy2vendorCmd())
	cmd.AddCommand(db.Cmds()...)
	return cmd
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
