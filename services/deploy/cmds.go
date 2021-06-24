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
			d.Build.Env, d.Push.Env = env, env
			if len(args) > 0 && len(args[0]) == 1 {
				d.Build.Tag, d.Push.Tag = args[0][0], args[0][0]
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
	cmd.Flags().BoolVarP(&d.alwaysPush, `always-push`, `P`, false,
		`Always push the built images to registry, even the deploying cluster is local machine.`)
	cmd.Flags().StringVarP(&d.DockerLogin.User, "docker-user", "u", "",
		"The docker login user, used before pushing or pulling image")
	cmd.Flags().StringVarP(&d.DockerLogin.Password, "docker-password", "p", "",
		"The docker login password, used before pushing or pulling image")

	cmd.Flags().StringVarP(&d.filter, `filter`, `f`, ``, `Filter the node to deploy to by node addr.`)

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
