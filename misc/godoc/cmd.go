package godoc

import (
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `godoc`,
		Short: `The godoc server on local machine.`,
	}
	cmd.AddCommand(runCmd())
	cmd.AddCommand(deployCmd())
	cmd.AddCommand(rmDeployCmd())
	cmd.AddCommand(accessCmd())
	cmd.AddCommand(shellCmd())
	cmd.AddCommand(psCmd())
	return cmd
}

func runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `run`,
		Short: `Run the godoc server using nohup.`,
		RunE:  release.NoArgCall(run),
	}
}

func deployCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `deploy`,
		Short: `Deploy the godoc server using docker image.`,
		RunE:  release.NoArgCall(deploy),
	}
}

func rmDeployCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `rm-deploy`,
		Short: `Stop and remove docker container of the godoc server.`,
		RunE:  release.NoArgCall(deploy),
	}
}

func accessCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `access`,
		Short: `Access config for the godoc server.`,
		RunE: release.NoArgCall(func() error {
			return accessPrint()
		}),
	}
	cmd.AddCommand(&cobra.Command{
		Use:   `setup`,
		Short: `Setup access config for the godoc server.`,
		RunE: release.NoArgCall(func() error {
			return accessSetup()
		}),
	})
	return cmd
}

func shellCmd() *cobra.Command {
	theCmd := &cobra.Command{
		Use:   `shell`,
		Short: `Enter the docker container of the godoc server.`,
		RunE:  release.NoArgCall(shell),
	}
	return theCmd
}

func psCmd() *cobra.Command {
	theCmd := &cobra.Command{
		Use:   `ps`,
		Short: `List the docker container of the godoc server.`,
		RunE:  release.NoArgCall(ps),
	}
	return theCmd
}
