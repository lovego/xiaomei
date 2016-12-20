package oam

import (
	"github.com/spf13/cobra"
)

func Cmds(serverFilter *string) []*cobra.Command {
	return []*cobra.Command{
		{
			Use: `restart`,
			Run: func(c *cobra.Command, args []string) {
				Restart(*serverFilter)
			},
		},
		{
			Use: `status`,
			Run: func(c *cobra.Command, args []string) {
				Status(*serverFilter)
			},
		},
		{
			Use: `shell`,
			Run: func(c *cobra.Command, args []string) {
				Shell(*serverFilter)
			},
		},
		{
			Use: `exec <cmd> [<args>...]`,
			Run: func(c *cobra.Command, args []string) {
				Exec(*serverFilter, args)
			},
		},
	}
}
