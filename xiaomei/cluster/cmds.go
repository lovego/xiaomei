package cluster

import (
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func LsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `ls`,
		Short: `list all nodes.`,
		RunE: release.EnvCall(func(env string) error {
			Get(env).List()
			return nil
		}),
	}
}
