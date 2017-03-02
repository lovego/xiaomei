package project

import (
	new "github.com/bughou-go/xiaomei/cli/project/new"
	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/bughou-go/xiaomei/cli/z"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `pj`,
		Short: `the whole project.`,
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   `deploy <env>`,
			Short: `deploy the project.`,
			RunE: z.Arg1Call(`dev`, func(env string) error {
				return stack.Deploy(env, ``)
			}),
		},
		&cobra.Command{
			Use:   `ps [<env>]`,
			Short: `list tasks of the project.`,
			RunE: z.Arg1Call(`dev`, func(env string) error {
				return stack.Ps(env, ``)
			}),
		},
		new.Cmd(),
	)
	return cmd
}
