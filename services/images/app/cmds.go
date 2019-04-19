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
	log.Println(color.GreenString(`compile the app binary.`))
	o := cmd.O{Dir: filepath.Join(release.Root(), `..`)}
	if linuxAMD64 && (runtime.GOOS != `linux` || runtime.GOARCH != `amd64`) {
		o.Env = []string{`GOOS=linux`, `GOARCH=amd64`} // cross compile
	}

	var subCmd []string
	if filepath.Base(o.Dir) == release.Name() && len(o.Env) == 0 {
		// go install have no "-o" option, so we use GOBIN environment variable
		// go install cannot install cross-compiled binaries when GOBIN is set
		o.Env = []string{`GOBIN=` + filepath.Join(o.Dir, `release/img-app`)}
		subCmd = []string{`install`, `-v`}
	} else {
		// after go 1.10, go build is as fast as go install.
		// so maybe we should always use go build in the future.
		subCmd = []string{`build`, `-i`, `-v`, `-o`, `release/img-app/` + release.Name()}
	}
	if cmd.Ok(o, `go`, subCmd...) {
		return misc.SpecAll()
	}
	return errors.New(`compile the app binary failed.`)
}
