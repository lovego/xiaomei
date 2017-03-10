package app

import (
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/spf13/cobra"
)

func DepsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `deps`,
		Short: `list all dependence packages.`,
		Run: func(c *cobra.Command, args []string) {
			listDeps()
		},
	}
	return cmd
}

func listDeps() {
	deps, _ := cmd.Run(cmd.O{Output: true}, `go`, `list`, `-e`, `-f`, `'{{join .Deps "\n" }}'`)
	pkgs := []string{}
	for _, dep := range strings.Split(deps, "\n") {
		if strings.Contains(dep, `.`) &&
			!strings.HasPrefix(dep, `github.com/bughou-go/xiaomei`) {
			pkgs = append(pkgs, dep)
		}
	}
	cmd.Run(cmd.O{}, `echo`, strings.Join(pkgs, "\n"))
}
