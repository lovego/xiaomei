package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func SpecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `spec`,
		Short: `check code spec.`,
		RunE: func(c *cobra.Command, args []string) error {
			arg := ``
			if len(args) > 0 {
				arg = args[0]
			}
			return Spec(arg)
		},
	}
}

func Spec(t string) error {
	config.Log(color.GreenString(`check app code spec.`))
	if err := os.Chdir(filepath.Join(config.Root(), `..`)); err != nil {
		panic(err)
	}

	targets := specTargets()
	if t != `all` {
		targets = specChangedTargets(targets)
	}
	if len(targets) == 0 {
		return nil
	}

	if !cmd.Ok(cmd.O{NoStdout: true}, `which`, `gospec`) {
		cmd.Run(cmd.O{Panic: true}, `go`, `get`, `-v`, `github.com/bughou-go/spec/gospec`)
	}

	if cmd.Ok(cmd.O{}, `gospec`, targets...) {
		return nil
	}
	return errors.New(`spec check failed.`)
}

func specTargets() []string {
	matches, err := filepath.Glob(`*`)
	if err != nil {
		panic(err)
	}

	targets := []string{}
	for _, v := range matches {
		if v != `release` && v != `vendor` {
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
