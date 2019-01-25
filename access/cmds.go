package access

import (
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

// access commands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `access [<env>]`,
		Short: `access config for the project.`,
		RunE: release.EnvCall(func(env string) error {
			return printNginxConf(env)
		}),
	}
	cmd.AddCommand(accessSetupCmd())
	cmd.AddCommand(accessReloadCmd())
	return cmd
}

func accessSetupCmd() *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `setup [<env>]`,
		Short: `setup access config for the project.`,
		RunE: release.EnvCall(func(env string) error {
			return SetupNginx(env, filter, "")
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	return cmd
}

func accessReloadCmd() *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `reload [<env>]`,
		Short: `reload access config for the project.`,
		RunE: release.EnvCall(func(env string) error {
			return ReloadNginx(env, filter)
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	return cmd
}
