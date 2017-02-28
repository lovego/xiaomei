package app

import (
	"fmt"

	"github.com/bughou-go/xiaomei/cli/cluster"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/spf13/cobra"
)

func PsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `ps [<env>]`,
		Short: `list tasks of app service.`,
		RunE: func(c *cobra.Command, args []string) error {
			env := `dev`
			if len(args) > 0 {
				env = args[0]
			}
			return Ps(env)
		},
	}
}

func Ps(env string) error {
	clusterConf, err := cluster.GetConfig(env)
	if err != nil {
		return err
	}
	addr, err2 := clusterConf.SshAddr()
	if err2 != nil {
		return err2
	}
	_, err = cmd.SshRun(cmd.O{}, addr,
		fmt.Sprintf(`docker service ps %s_app`, config.DeployName()),
	)
	return err
}

func Restart(env string) error {
	return nil
}

func Shell(env string) error {
	return nil
}

func Exec(env string) error {
	return nil
}
