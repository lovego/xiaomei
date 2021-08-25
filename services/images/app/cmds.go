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
	"github.com/lovego/config/config"
	"github.com/lovego/xiaomei/misc"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		runCmd(`run`),
		compileCmd(),
		depsCmd(),
		copy2vendorCmd(),
		runCmd(`doc`),
	}
}

func runCmd(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   name + ` [flags] [env] [-- go build flags]`,
		DisableFlagsInUseLine: true,
		RunE: release.EnvSliceCall(func(env string, args []string) error {
			if err := Compile(false, env, args); err != nil {
				return err
			}
			signal.Ignore(os.Interrupt)
			o := cmd.O{
				Dir: filepath.Join(release.Root(), `..`),
				Env: []string{`ProDEV=true`, config.EnvVar + `=` + env},
			}
			if name == `doc` {
				o.Env = append(o.Env, `GOA_DOC=1`)
			}
			_, err := cmd.Run(o, filepath.Join(release.ServiceDir(`app`), release.Name(env)))
			return err
		}),
	}
	if name == `doc` {
		cmd.Short = `Compile the app server binary and run it to generate routes documents.`
	} else {
		cmd.Short = `Compile the app server binary and run it.`
	}
	return cmd
}

func compileCmd() *cobra.Command {
	return &cobra.Command{
		Use:                   `compile [flags] [env] [-- go build flags]`,
		DisableFlagsInUseLine: true,
		Short:                 `Compile the app server binary.`,
		RunE: release.EnvSliceCall(func(env string, args []string) error {
			return Compile(false, env, args)
		}),
	}
}

func Compile(linuxAMD64 bool, env string, goBuildFlags []string) error {
	log.Println(color.GreenString(`compile the app binary.`))
	o := cmd.O{Dir: release.SrcDir()}
	if linuxAMD64 && (runtime.GOOS != `linux` || runtime.GOARCH != `amd64`) {
		o.Env = []string{`GOOS=linux`, `GOARCH=amd64`} // cross compile
	}
	var options = []string{
		`build`, `-v`, `-o`, filepath.Join(release.ServiceDir(`app`), release.Name(env)),
	}
	options = append(options, goBuildFlags...)
	if cmd.Ok(o, release.GoCmd(), options...) {
		return misc.SpecAll()
	}
	return errors.New(`compile the app binary failed.`)
}
