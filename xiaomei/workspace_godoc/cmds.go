package workspace_godoc

import (
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `workspace-godoc`,
		Short: `deploy the workspace godoc server.`,
		RunE:  release.NoArgCall(deploy),
	}
	cmd.AddCommand(accessCmd())
	cmd.AddCommand(shellCmd())
	return cmd
}

func accessCmd() *cobra.Command {
	var setup bool
	cmd := &cobra.Command{
		Use:   `access [<env>]`,
		Short: `access config for the workspace godoc server.`,
		RunE: release.NoArgCall(func() error {
			if setup {
				return accessSetup()
			} else {
				return accessPrint()
			}
		}),
	}
	cmd.Flags().BoolVarP(&setup, `setup`, `s`, false, `setup access.`)
	return cmd
}

func shellCmd() *cobra.Command {
	theCmd := &cobra.Command{
		Use:   `shell [<env>]`,
		Short: `enter the container for workspace godoc server.`,
		RunE:  release.NoArgCall(shell),
	}
	return theCmd
}
