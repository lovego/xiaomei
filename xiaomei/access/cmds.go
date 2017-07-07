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
	cmd := &cobra.Command{
		Use:   `access`,
		Short: `access config for the ` + desc(svcName) + `.`,
		RunE: release.NoArgCall(func() error {
			if setup {
				return accessSetup(svcName)
			} else {
				return accessPrint(svcName)
			}
		}),
	}
	cmd.Flags().BoolVarP(&setup, `setup`, `s`, false, `setup access.`)
	return cmd
}

func desc(svcName string) string {
	if svcName == `` {
		return `project`
	} else {
		return svcName + ` service`
	}
}
