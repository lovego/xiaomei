package access

import (
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

// Run, Deploy, Ps, Logs commands
func Cmds(svcName string) (cmds []*cobra.Command) {
	if svcName == `` || svcName == `app` || svcName == `web` || svcName == `godoc` {
		cmds = append(cmds, accessCmd(svcName))
	}
	return
}

func accessCmd(svcName string) *cobra.Command {
	var setup bool
	var filter string
	cmd := &cobra.Command{
		Use:   `access [<env>]`,
		Short: `access config for the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			if setup {
				return accessSetup(env, svcName, filter)
			} else {
				return accessPrint(env, svcName)
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
