package tasks

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `compile`,
			Short: `compile the tasks binary.`,
			RunE:  release.NoArgCall(compile),
		},
	}
}

func compile() error {
	log.Println(color.GreenString(`compile the tasks binary.`))
	if cmd.Ok(cmd.O{
		Dir: filepath.Join(release.Root(), "../tasks"),
		Env: []string{
			`GOBIN=` + filepath.Join(release.Root(), `img-app`),
		},
	}, `go`, `install`, `-v`) {
		return nil
	}
	return errors.New(`compile the tasks binary failed.`)
}
