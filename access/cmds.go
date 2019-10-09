package access

import (
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

const sudoTip = `
The user must be permitted to run some commands with sudo. A line like this in /etc/sudoers may work:
  USERNAME ALL=NOPASSWD: /bin/tee, /bin/mkdir, /usr/sbin/nginx -t, /bin/systemctl reload nginx`

// access commands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `access [<env>]`,
		Short: `Access config for the project.`,
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
		Use:   `setup [<env>] [flags]` + sudoTip,
		Short: `Setup access config for the project.`,
		DisableFlagsInUseLine: true,
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
		Use:   `reload [<env>] [flags]` + sudoTip,
		Short: `Reload access config for the project.`,
		DisableFlagsInUseLine: true,
		RunE: release.EnvCall(func(env string) error {
			return ReloadNginx(env, filter)
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	return cmd
}
