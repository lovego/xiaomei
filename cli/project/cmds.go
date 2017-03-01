package project

import (
	"errors"

	new "github.com/bughou-go/xiaomei/cli/project/new"
	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `pj`,
		Short: `the whole project.`,
	}
	cmd.AddCommand(
		deployCmd(),
		psCmd(),
		new.Cmd(),
	)
	return cmd
}

func deployCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `deploy <env>`,
		Short: `deploy project to the specified env.`,
		RunE: func(c *cobra.Command, args []string) error {
			var env string
			if len(args) > 1 {
				return errors.New(`redundant args.`)
			} else if len(args) == 1 {
				env = args[0]
			}
			return stack.Deploy(env, ``)
		},
	}
}
