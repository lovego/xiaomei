package project

import (
	"github.com/bughou-go/xiaomei/cli/app"
	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/bughou-go/xiaomei/cli/z"
	"github.com/spf13/cobra"
)

func buildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `deploy <env>`,
		Short: `deploy the project.`,
		RunE: z.Arg1Call(`dev`, func(env string) error {
			return build(env, ``)
		}),
	}
}

func build() error {
	stack.Services()
}
