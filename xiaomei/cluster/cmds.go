package cluster

import (
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func LsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `ls`,
		Short: `list all clusters.`,
		RunE: release.NoArgCall(func() error {
			GetCluster().List()
			return nil
		}),
	}
}
