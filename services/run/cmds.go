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
	var pull bool
	cmd := &cobra.Command{
		Use:   `run [<env>]`,
		Short: `run    the ` + svcName + ` service.`,
		RunE: release.EnvCall(func(env string) error {
			if err := images.Build(svcName, env, ``, pull); err != nil {
				return err
			}
			return run(env, svcName)
		}),
	}
	cmd.Flags().BoolVarP(&pull, `pull`, `p`, true, `pull base image.`)
	return cmd
}
