package tasks

import (
	"errors"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `build-bin`,
			Short: `build the tasks.`,
			RunE:  release.NoArgCall(buildBinary),
		},
	}
}

func buildBinary() error {
	utils.Log(color.GreenString(`building tasks binary.`))
	if cmd.Ok(cmd.O{
		Dir: filepath.Join(release.Root(), "../tasks"),
		Env: []string{`GOBIN=` + filepath.Join(release.Root(), `img-app`)},
	}, `go`, `install`, `-v`) {
		return nil
	}
	return errors.New(`building tasks binary failed.`)
}
