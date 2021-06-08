package run

import (
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services/images"
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
		Use:   `run [env]`,
		Short: `Run    the ` + svcName + ` service.`,
		RunE: release.EnvCall(func(env string) error {
			if err := (images.Build{Env: env}).Run(svcName); err != nil {
				return err
			}
			return run(env, svcName)
		}),
	}
	return cmd
}
