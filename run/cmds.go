package run

import (
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds(svcName string) (cmds []*cobra.Command) {
	if svcName != `` {
		cmds = append(cmds, runCmdFor(svcName))
	}
	return
}

func runCmdFor(svcName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `run [<env>]`,
		Short: `run    the ` + svcName + ` service.`,
		RunE: release.EnvCall(func(env string) error {
			return run(env, svcName)
		}),
	}
	return cmd
}
