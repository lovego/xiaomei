package deploy

import (
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/images/registry"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

// Run, Deploy, Ps, Logs commands
func Cmds(svcName string) (cmds []*cobra.Command) {
	if svcName != `` {
		cmds = append(cmds, runCmdFor(svcName))
	}
	cmds = append(cmds,
		deployCmdFor(svcName),
		rmDeployCmdFor(svcName),
	)
	return
}

func runCmdFor(svcName string) *cobra.Command {
	// var publish []string
	cmd := &cobra.Command{
		Use:   `run [<env>]`,
		Short: `run    the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return run(env, svcName)
		}),
	}
	// cmd.Flags().StringSliceVarP(&publish, `publish`, `p`, nil, `publish ports for container.`)
	return cmd
}

func deployCmdFor(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `deploy [<env> [<tag>]]`,
		Short: `deploy the ` + desc(svcName) + `.`,
		RunE: release.Env1Call(func(env, timeTag string) error {
			noTag := timeTag == ``
			if timeTag == `` {
				timeTag = conf.TimeTag(env)
				if err := images.Build(svcName, env, timeTag, true); err != nil {
					return err
				}
				if err := images.Push(svcName, env, timeTag); err != nil {
					return err
				}
			}
			if err := deploy(svcName, env, timeTag, filter); err != nil {
				return err
			}
			if noTag {
				registry.PruneTimeTags(svcName, env, 10)
			}
			return nil
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	return cmd
}

func rmDeployCmdFor(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `rm-deploy [<env>]`,
		Short: `remove deployment of the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return rmDeploy(svcName, env, filter)
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
