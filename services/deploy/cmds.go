package deploy

import (
	"errors"

	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds(svcName string) (cmds []*cobra.Command) {
	return []*cobra.Command{
		deployCmdFor(svcName),
		rmDeployCmdFor(svcName),
	}
}

func deployCmdFor(svcName string) *cobra.Command {
	var d = Deploy{svcName: svcName}
	cmd := &cobra.Command{
		Use:   `deploy [flags] [env [tag]] [ -- [prepare flags] [-- docker build flags] ]`,
		Short: `Deploy the ` + desc(svcName) + `.`,
		RunE: release.EnvSlicesCall(func(env string, args [][]string) error {
			if len(args) > 3 || len(args) > 0 && len(args[0]) > 1 {
				return errors.New("invalid arguments usage.")
			}
			d.Env = env
			if len(args) > 0 && len(args[0]) == 1 {
				d.Tag = args[0][0]
			}
			if len(args) > 1 {
				d.PrepareFlags = args[1]
			}
			if len(args) > 2 {
				d.DockerBuildFlags = args[2]
			}
			return d.start()
		}),
		DisableFlagsInUseLine: true,
	}
	cmd.Flags().BoolVarP(&d.noPushImageIfLocal, `no-push-if-local`, `P`, false,
		`Don't push the built images to registry if deploying cluster is local machine.`)

	cmd.Flags().StringVarP(&d.filter, `filter`, `f`, ``, `Filter the node to deploy to by node addr.`)

	cmd.Flags().StringVarP(&d.beforeScript, `before-script`, `b`, ``,
		`Script to be executed on:
  1. This local machine at the very beginning of the deployment.
  2. Every node in the cluster before deployment on the node.
If this local machine is also a node in the cluster, and the
before-script is aready executed by the "1" step, it won't be
executed by the "2" step again.`)

	cmd.Flags().BoolVarP(&d.noBeforeScriptOnLocal, `no-before-script-on-local`, `B`, false,
		`Don't run before-script on this local machine at the very beginning of the deployment(the "1" step).`)

	cmd.Flags().BoolVarP(&d.noWatch, `no-watch`, `W`, false,
		`After deployed a node, don't watch container status until "Ctl+C".`)

	return cmd
}

func rmDeployCmdFor(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `rm-deploy [env]`,
		Short: `Remove deployment of the ` + desc(svcName) + `.`,
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
