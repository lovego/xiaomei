package cluster

import (
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `cluster`,
		Short: `ssh into a manager of the cluster.`,
		RunE: release.NoArgCall(func() error {
			_, err := Run(cmd.O{}, ``)
			return err
		}),
	}
	cmd.AddCommand(lsCmd())
	return cmd
}

func lsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `ls`,
		Short: `list all clusters.`,
		RunE: release.NoArgCall(func() error {
			GetCluster().List()
			return nil
		}),
	}
}
