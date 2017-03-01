package project

import (
	"fmt"

	"github.com/bughou-go/xiaomei/cli/cluster"
	"github.com/bughou-go/xiaomei/config"
	"github.com/spf13/cobra"
)

func psCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `ps [<env>]`,
		Short: `list tasks of app service.`,
		RunE: func(c *cobra.Command, args []string) error {
			env := `dev`
			if len(args) > 0 {
				env = args[0]
			}
			return ps(env)
		},
	}
}

func ps(env string) error {
	return cluster.Run(env, fmt.Sprintf(`docker stack ps %s`, config.DeployName()))
}

func restart(env string) error {
	return nil
}

func shell(env string) error {
	return nil
}

func exec(env string) error {
	return nil
}
