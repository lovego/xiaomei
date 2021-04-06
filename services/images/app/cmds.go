package app

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/misc"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		runCmd(),
		{
			Use:   `compile`,
			Short: `Compile the app server binary.`,
			RunE: release.NoArgCall(func() error {
				return compile(false)
			}),
		},
		depsCmd(),
		copy2vendorCmd(),
	}
}

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `run [<env>]`,
		Short: `Compile and run the app server binary.`,
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
	log.Println(color.GreenString(`compile the app binary.`))
	o := cmd.O{Dir: filepath.Join(release.Root(), `..`)}
	if linuxAMD64 && (runtime.GOOS != `linux` || runtime.GOARCH != `amd64`) {
		o.Env = []string{`GOOS=linux`, `GOARCH=amd64`} // cross compile
	}
	if cmd.Ok(o, release.GoCmd(), `build`, `-i`, `-v`, `-o`, `release/img-app/`+release.Name()) {
		return misc.SpecAll()
	}
	return errors.New(`compile the app binary failed.`)
}
