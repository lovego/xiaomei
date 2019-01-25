package godoc

import (
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `godoc`,
		Short: `[service] the workspace godoc server.`,
	}
	cmd.AddCommand(deployCmd())
	cmd.AddCommand(accessCmd())
	cmd.AddCommand(shellCmd())
	cmd.AddCommand(psCmd())
	return cmd
}

func deployCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `deploy`,
		Short: `deploy the workspace godoc server.`,
		RunE:  release.NoArgCall(deploy),
	}
}

func accessCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `access`,
		Short: `access config for the workspace godoc server.`,
		RunE: release.NoArgCall(func() error {
			return accessPrint()
		}),
	}
	cmd.AddCommand(&cobra.Command{
		Use:   `setup`,
		Short: `setup access config for the workspace godoc server.`,
		RunE: release.NoArgCall(func() error {
			return accessSetup()
		}),
	})
	return cmd
}

func shellCmd() *cobra.Command {
	theCmd := &cobra.Command{
		Use:   `shell`,
		Short: `enter the container for workspace godoc server.`,
		RunE:  release.NoArgCall(shell),
	}
	return theCmd
}

func psCmd() *cobra.Command {
	theCmd := &cobra.Command{
		Use:   `ps`,
		Short: `list the container for workspace godoc server.`,
		RunE:  release.NoArgCall(ps),
	}
	return theCmd
}
