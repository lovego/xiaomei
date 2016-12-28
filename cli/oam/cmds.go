package oam

import (
	"github.com/spf13/cobra"
)

func Cmds(serverFilter *string) []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `restart`,
			Short: `restart the specified environment appserver`,
			Run: func(c *cobra.Command, args []string) {
				Restart(*serverFilter)
			},
		},
		{
			Use:   `status`,
			Short: `check the specified environment appserver status`,
			Run: func(c *cobra.Command, args []string) {
				Status(*serverFilter)
			},
		},
		{
			Use:   `shell`,
			Short: `into the specified environment command line`,
			Run: func(c *cobra.Command, args []string) {
				Shell(*serverFilter)
			},
		},
		{
			Use:   `exec <cmd> [<args>...]`,
			Short: `execute the specified environment command`,
			Run: func(c *cobra.Command, args []string) {
				Exec(*serverFilter, args)
			},
		},
	}
}
