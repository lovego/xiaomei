package app

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/fs"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func specCmd() *cobra.Command {
	var onlyChanged bool
	cmd := &cobra.Command{
		Use:   `spec`,
		Short: `check code spec.`,
		RunE: func(c *cobra.Command, args []string) error {
			return gospec(args, onlyChanged)
		},
	}
	cmd.Flags().BoolVarP(&onlyChanged, `only-changed`, `c`, true, `only check the changed files.`)
	return cmd
}

func gospec(targets []string, onlyChanged bool) error {
	log.Println(color.GreenString(`check app code spec.`))
	if err := os.Chdir(filepath.Join(release.Root(), `..`)); err != nil {
		panic(err)
	}

	if len(targets) == 0 {
		targets = specTargets()
	}
	if onlyChanged {
		targets = specChangedTargets(targets)
	}
	if len(targets) == 0 {
		return nil
	}

	if !cmd.Ok(cmd.O{NoStdout: true}, `which`, `gospec`) {
		cmd.Run(cmd.O{Panic: true}, `go`, `get`, `-v`, `github.com/lovego/spec/gospec`)
	}

	if cmd.Ok(cmd.O{}, `gospec`, targets...) {
		return nil
	}
	return errors.New(`spec check failed.`)
}

func specTargets() []string {
	matches, err := filepath.Glob(`[^.]*`)
	if err != nil {
		panic(err)
	}

	targets := []string{}
	for _, v := range matches {
		if v != `release` && v != `vendor` {
			if fs.IsDir(v) {
				v = v + `/`
			}
			targets = append(targets, v)
		}
	}
	return targets
}

func specChangedTargets(targets []string) []string {
	output, _ := cmd.Run(cmd.O{Output: true}, `git`, append(
		[]string{`diff`, `--name-only`, `--diff-filter=AMR`, `--relative`, `--`}, targets...,
	)...)
	lines := strings.Split(output, "\n")
	results := []string{}
	for _, v := range lines {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			results = append(results, v)
		}
	}
	return results
}
