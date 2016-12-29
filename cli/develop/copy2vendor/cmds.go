package copy2vendor

import (
	"errors"
	"github.com/spf13/cobra"
)

func Cmds() *cobra.Command {
	var n bool
	copy := &cobra.Command{
		Use:   `copy2vendor <package>...`,
		Short: `copy the specified packages to project vendor dir.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) <= 0 {
				return errors.New(`need at least a package path`)
			}
			return Copy2Vendor(args, n)
		},
	}
	flags := copy.Flags()
	flags.BoolVarP(&n, `no-clobber`, `n`, false, `do not overwrite an existing file.`)
	return copy
}
