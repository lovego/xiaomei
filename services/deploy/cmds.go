package deploy

import (
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
	var d Deploy
	cmd := &cobra.Command{
		Use:   `deploy [<env> [<tag>]]`,
		Short: `Deploy the ` + desc(svcName) + `.`,
		RunE: release.Env1Call(func(env, timeTag string) error {
			d.svcName, d.env, d.timeTag = svcName, env, timeTag
			return d.start()
		}),
	}
	cmd.Flags().BoolVar(&d.noPullBaseImage, `no-pull`, false,
		`Don't pull base image when building image.`)
	cmd.Flags().BoolVar(&d.noPushImageIfLocal, `no-push`, false,
		`Don't push the built images to registry if deploying cluster is local machine.`)
	cmd.Flags().StringVarP(&d.filter, `filter`, `f`, ``, `Filter the node to deploy to by node addr.`)

	cmd.Flags().StringVarP(&d.beforeScript, `before-script`, `b`, ``,
		`Script to execute before deploy on every node.`)
	cmd.Flags().BoolVarP(&d.noBeforeScriptOnLocal, `no-before-script-on-local`, `B`, false,
		`Don't run before-script on this local machine at the very beginning of the deployment.`)
	cmd.Flags().BoolVarP(&d.noWatch, `no-watch`, `W`, false,
		`After deployed a node, don't watch container status until "Ctl+C".`)
	return cmd
}

func rmDeployCmdFor(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `rm-deploy [<env>]`,
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
