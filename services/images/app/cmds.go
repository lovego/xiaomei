package app

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/misc"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		execCmd(),
		{
			Use:   `compile`,
			Short: `compile the app server binary.`,
			RunE: release.NoArgCall(func() error {
				return compile(false)
			}),
		},
		depsCmd(),
		copy2vendorCmd(),
	}
}

func execCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `exec [<env>]`,
		Short: `compile and execute the app server binary.`,
		RunE: release.EnvCall(func(env string) error {
			if err := compile(false); err != nil {
				return err
			}
			signal.Ignore(os.Interrupt)
			_, err := cmd.Run(cmd.O{
				Dir: filepath.Join(release.Root(), `..`),
				Env: []string{
					`GODEV=true`, Image{}.EnvironmentEnvVar() + `=` + env,
				},
			}, filepath.Join(release.Root(), `img-app`, release.Name()))
			return err
		}),
	}
	return cmd
}

func compile(linuxAMD64 bool) error {
	log.Println(color.GreenString(`compile the app server binary.`))
	o := cmd.O{Dir: filepath.Join(release.Root(), `..`)}
	if linuxAMD64 {
		o.Env = []string{"GOOS=linux", "GOARCH=amd64"}
	}
	if cmd.Ok(o, `go`, `build`, `-v`, `-o`, `release/img-app/`+release.Name()) {
		return misc.SpecAll()
	}
	return errors.New(`compile the app server binary failed.`)
}
