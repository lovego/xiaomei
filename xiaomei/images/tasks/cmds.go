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

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `tasks`,
		Short: `the tasks.`,
	}
	cmd.AddCommand(&cobra.Command{
		Use:   `build-bin`,
		Short: `build the tasks.`,
		RunE:  release.NoArgCall(buildBinary),
	})
	return cmd
}

func buildBinary() error {
	utils.Log(color.GreenString(`building tasks binary.`))
	if cmd.Ok(cmd.O{
		Dir: filepath.Join(release.Root(), "../tasks"),
		Env: []string{`GOBIN=` + filepath.Join(release.Root(), `img-app`)},
	}, `go`, `install`) {
		return nil
	}
	return errors.New(`building tasks binary failed.`)
}
