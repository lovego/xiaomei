package oam

import (
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		makeCmd(`status`, `show appserver status`, Status),
		makeCmd(`restart`, `restart appserver`, Restart),
		makeCmd(`shell`, `enter shell`, Shell),
		makeCmd(`exec <cmd> [<args>...]`, `execute command`, Exec),
	}
}

func makeCmd(use, short string, fun func(filter string, args []string) error) *cobra.Command {
	var filter *string
	cmd := &cobra.Command{
		Use:   use,
		Short: `[oam] ` + short,
		RunE: func(c *cobra.Command, args []string) error {
			return fun(*filter, args)
		},
	}
	filter = cmd.Flags().StringP(`server`, `s`, ``, `match servers by Addr or Tasks.`)
	return cmd
}
