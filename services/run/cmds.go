package run

import (
	"errors"

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
		Use:   `run [env] [ -- [prepare flags] [-- docker build flags] ]`,
		Short: `Run    the ` + svcName + ` service.`,
		RunE: release.EnvSlicesCall(func(env string, args [][]string) error {
			b := images.Build{Env: env}
			if len(args) > 0 && len(args[0]) > 0 {
				return errors.New("invalid arguments usage.")
			}
			if len(args) > 1 {
				b.PrepareFlags = args[2]
			}
			if len(args) > 2 {
				b.DockerBuildFlags = args[2]
			}
			if err := b.Run(svcName); err != nil {
				return err
			}
			return run(env, svcName)
		}),
	}
	return cmd
}
