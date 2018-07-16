package access

import (
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

// access commands
func Cmd(svcName string) *cobra.Command {
	if svcName == `` || svcName == `app` || svcName == `web` || svcName == `godoc` {
		return accessCmd(svcName)
	}
	return nil
}

func accessCmd(svcName string) *cobra.Command {
	var setup bool
	var filter string
	cmd := &cobra.Command{
		Use:   `access [<env>]`,
		Short: `access config for the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			if setup {
				return SetupNginx(env, svcName, filter, "")
			} else {
				return printNginxConf(env, svcName)
			}
		}),
	}
	cmd.Flags().BoolVarP(&setup, `setup`, `s`, false, `setup access.`)
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
