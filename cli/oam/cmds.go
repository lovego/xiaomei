package oam

import (
	"github.com/spf13/cobra"
)

func Cmds(serverFilter *string) []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `restart`,
			Short: `[oam] restart appserver`,
			Run: func(c *cobra.Command, args []string) {
				Restart(*serverFilter)
			},
		},
		{
			Use:   `status`,
			Short: `[oam] show appserver status`,
			Run: func(c *cobra.Command, args []string) {
				Status(*serverFilter)
			},
		},
		{
			Use:   `shell`,
			Short: `[oam] enter shell`,
			Run: func(c *cobra.Command, args []string) {
				Shell(*serverFilter)
			},
		},
		{
			Use:   `exec <cmd> [<args>...]`,
			Short: `[oam] execute command`,
			Run: func(c *cobra.Command, args []string) {
				Exec(*serverFilter, args)
			},
		},
	}
}
