package app

import (
	"github.com/bughou-go/xiaomei/cli/app/deps"
	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/bughou-go/xiaomei/cli/z"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `app`,
		Short: `the appserver.`,
	}
	cmd.AddCommand(
		RunCmd(),
		BuildCmd(),
		&cobra.Command{
			Use:   `deploy [<env>]`,
			Short: `deploy the app service.`,
			RunE: z.Arg1Call(`env`, func(env string) error {
				return stack.Ps(env, `app`)
			}),
		},
		&cobra.Command{
			Use:   `ps [<env>]`,
			Short: `list tasks of app service.`,
			RunE: z.Arg1Call(`env`, func(env string) error {
				return stack.Ps(env, `app`)
			}),
		},
		deps.Cmd(),
		SpecCmd(),
	)
	return cmd
}
