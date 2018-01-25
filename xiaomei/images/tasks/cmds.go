package tasks

import (
	"errors"
	"log"
	"os"
	"os/signal"
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
		execCmd(),
	}
}

func execCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `exec [<env>]`,
		Short: `compile and execute the app server binary.`,
		RunE: release.EnvCall(func(env string) error {
			if err := compile(); err != nil {
				return err
			}
			signal.Ignore(os.Interrupt)
			_, err := cmd.Run(cmd.O{
				Dir: filepath.Join(release.Root(), `..`),
				Env: []string{
					`GODEV=true`, Image{}.EnvironmentEnvName() + `=` + env,
				},
			}, filepath.Join(release.Root(), `img-app/tasks`))
			return err
		}),
	}
	return cmd
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
