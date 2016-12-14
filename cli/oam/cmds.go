package oam

import (
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use: `restart`,
			Run: func(c *cobra.Command, args []string) {
				Restart()
			},
		},
		{
			Use: `status`,
			Run: func(c *cobra.Command, args []string) {
				Status()
			},
		},
		{
			Use: `shell`,
			Run: func(c *cobra.Command, args []string) {
				Shell()
			},
		},
		{
			Use: `exec <cmd> [<args>...]`,
			Run: func(c *cobra.Command, args []string) {
				Exec(args)
			},
		},
	}
}
