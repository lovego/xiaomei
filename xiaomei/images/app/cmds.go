package app

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
		execCmd(),
		{
			Use:   `compile`,
			Short: `compile the app server binary.`,
			RunE:  release.NoArgCall(compile),
		},
		specCmd(),
		depsCmd(),
		copy2vendorCmd(),
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
			_, err := cmd.Run(cmd.O{
				Dir: filepath.Join(release.Root(), `..`),
				Env: []string{
					`GODEV=true`, Image{}.EnvironmentEnvName() + `=` + env,
				},
			}, filepath.Join(release.Root(), `img-app`, release.Name()))
			return err
		}),
	}
	return cmd
}

func compile() error {
	log.Println(color.GreenString(`compile the app server binary.`))
	if cmd.Ok(cmd.O{
		Dir: filepath.Join(release.Root(), `..`),
		Env: []string{
			`GOBIN=` + filepath.Join(release.Root(), `img-app`),
		},
	}, `go`, `install`, `-v`) {
		return gospec(nil, true)
	}
	return errors.New(`compile the app server binary failed.`)
}
