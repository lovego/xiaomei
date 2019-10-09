package deploy

import (
	"github.com/lovego/xiaomei/release"
	"github.com/lovego/xiaomei/services/images"
	"github.com/spf13/cobra"
)

func Cmds(svcName string) (cmds []*cobra.Command) {
	return []*cobra.Command{
		deployCmdFor(svcName),
		rmDeployCmdFor(svcName),
	}
}

func deployCmdFor(svcName string) *cobra.Command {
	var filter string
	var pull, push, noWatch bool
	cmd := &cobra.Command{
		Use:   `deploy [<env> [<tag>]]`,
		Short: `deploy the ` + desc(svcName) + `.`,
		RunE: release.Env1Call(func(env, timeTag string) error {
			if timeTag == `` {
				timeTag = release.TimeTag(env)
				if err := images.Build(svcName, env, timeTag, pull); err != nil {
					return err
				}
				if push {
					if err := images.Push(svcName, env, timeTag); err != nil {
						return err
					}
				}
			}
			if err := deploy(svcName, env, timeTag, filter, noWatch); err != nil {
				return err
			}
			return nil
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	cmd.Flags().BoolVarP(&pull, `pull`, `p`, true, `pull base image.`)
	cmd.Flags().BoolVarP(&push, `push`, `P`, true, `push the built images to registry.`)
	cmd.Flags().BoolVarP(&noWatch, `no-watch`, `W`, false, `after deployed a node, don't watch container status until "Ctl+C".`)
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
