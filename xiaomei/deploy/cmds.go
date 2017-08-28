package deploy

import (
	"github.com/lovego/xiaomei/xiaomei/images"
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
		psCmdFor(svcName),
		logsCmdFor(svcName),
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
	var noBuild, noPush /*, rmCurrent*/ bool
	var filter string
	cmd := &cobra.Command{
		Use:   `deploy [<env>]`,
		Short: `deploy the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			if !noBuild {
				if err := images.Build(env, svcName, true); err != nil {
					return err
				}
			}
			if !noPush {
				if err := images.Push(env, svcName); err != nil {
					return err
				}
			}
			return deploy(env, svcName, filter) // , rmCurrent)
		}),
	}
	cmd.Flags().BoolVarP(&noBuild, `no-build`, `B`, false, `do not build the images.`)
	cmd.Flags().BoolVarP(&noPush, `no-push`, `P`, false, `do not push the images.`)
	// cmd.Flags().BoolVar(&rmCurrent, `rm-current`, false, `remove the current running `+desc+`.`)
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	return cmd
}

func rmDeployCmdFor(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `rm-deploy [<env>]`,
		Short: `remove deployment of the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return rmDeploy(env, svcName, filter)
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	return cmd
}

func psCmdFor(svcName string) *cobra.Command {
	var watch bool
	var filter string
	cmd := &cobra.Command{
		Use:   `ps [<env>]`,
		Short: `list tasks of the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return ps(env, svcName, filter, watch)
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	cmd.Flags().BoolVarP(&watch, `watch`, `w`, false, `watch ps.`)
	return cmd
}

func logsCmdFor(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `logs [<env>]`,
		Short: `list logs  of the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return logs(env, svcName, filter)
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
