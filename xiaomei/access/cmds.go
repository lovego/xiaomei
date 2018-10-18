package access

import (
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

// access commands
func Cmd(svcName string) *cobra.Command {
	if svcName == `` || svcName == `godoc` {
		return accessCmd(svcName)
	}
	return nil
}

func accessCmd(svcName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `access [<env>]`,
		Short: `access config for the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return printNginxConf(env, svcName)
		}),
	}
	cmd.AddCommand(accessSetupCmd(svcName))
	if svcName == "" {
		cmd.AddCommand(accessReloadCmd(svcName))
	}
	return cmd
}

func accessSetupCmd(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `setup [<env>]`,
		Short: `setup access config for the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return SetupNginx(env, svcName, filter, "")
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	return cmd
}

func accessReloadCmd(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `reload [<env>]`,
		Short: `reload access config for the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return ReloadNginx(env, filter)
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	return cmd
}

func desc(svcName string) string {
	if svcName == `` {
		return `project`
	} else {
		return svcName + ` service`
	}
}
