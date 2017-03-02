package project

import (
	"github.com/bughou-go/xiaomei/cli/project/new"
	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/bughou-go/xiaomei/cli/z"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `pj`,
		Short: `the whole project.`,
	}
	cmd.AddCommand(&cobra.Command{
		Use:   `new <project-path>`,
		Short: `create a new project.`,
		RunE: z.Arg1Call(``, func(dir string) error {
			return new.New(dir)
		}),
	})
	cmd.AddCommand(stack.BDPcmds(``)...)
	return cmd
}
