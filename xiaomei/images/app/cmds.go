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

func compile() error {
	log.Println(color.GreenString(`compile the app server binary.`))
	if cmd.Ok(cmd.O{
		Dir: filepath.Join(release.Root(), `..`),
		Env: []string{`GOOS=linux`, `GOARCH=amd64`},
	}, `go`, `build`, `-i`, `-v`, `-o`, filepath.Join(release.Root(), `img-app`, release.Name())) {
		return gospec(nil, true)
	}
	return errors.New(`compile the app server binary failed.`)
}
