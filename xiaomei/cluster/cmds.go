package cluster

import (
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/z"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `cluster`,
		Short: `cluster operations.`,
	}
	cmd.AddCommand(lsCmd(), shellCmd())
	return cmd
}

func lsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `ls`,
		Short: `list all clusters.`,
		RunE: z.NoArgCall(func() error {
			GetCluster().List()
			return nil
		}),
	}
}

func shellCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `shell`,
		Short: `ssh into a manager of cluster.`,
		RunE: z.NoArgCall(func() error {
			_, err := Run(cmd.O{}, ``)
			return err
		}),
	}
}
