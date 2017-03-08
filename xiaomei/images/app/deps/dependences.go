package deps

import (
	"errors"
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `deps`,
		Short: `list all dependence packages.`,
		Run: func(c *cobra.Command, args []string) {
			List()
		},
	}
	cmd.AddCommand(copy2vendorCmd())
	return cmd
}

func copy2vendorCmd() *cobra.Command {
	var n bool
	cmd := &cobra.Command{
		Use:   `copy2vendor <package-path> ...`,
		Short: `copy the specified packages to vendor dir.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) <= 0 {
				return errors.New(`need at least one package path.`)
			}
			return Copy2Vendor(args, n)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&n, `no-clobber`, `n`, false, `do not overwrite an existing file.`)
	return cmd
}

func List() {
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
