package cluster

import (
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		shellCmd(),
		{
			Use:   `ls`,
			Short: `list all nodes.`,
			RunE: release.EnvCall(func(env string) error {
				Get(env).List()
				return nil
			}),
		},
	}
}

func shellCmd() *cobra.Command {
	var filter string
	theCmd := &cobra.Command{
		Use:   `shell [<env>]`,
		Short: `enter a node shell.`,
		RunE: release.EnvCall(func(env string) error {
			_, err := Get(env).Run(filter, cmd.O{}, ``)
			return err
		}),
	}
	theCmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr`)
	return theCmd
}
