package deploy

import (
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/registry"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

// Run, Deploy, Ps, Logs commands
func Cmds(svcName string) (cmds []*cobra.Command) {
	if svcName != `` {
		cmds = append(cmds, runCmdFor(svcName), shellCmdFor(svcName), tagsCmdFor(svcName))
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

func shellCmdFor(svcName string) *cobra.Command {
	var filter string
	theCmd := &cobra.Command{
		Use:   `shell [<env>]`,
		Short: `enter a container for ` + desc(svcName),
		RunE: release.EnvCall(func(env string) error {
			return shell(svcName, env, filter)
		}),
	}
	theCmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr`)
	return theCmd
}

func tagsCmdFor(svcName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `tags [<env>]`,
		Short: `list image tags of the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			registry.ListTimeTags(svcName, env)
			return nil
		}),
	}
	return cmd
}

func deployCmdFor(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `deploy [<env>]`,
		Short: `deploy the ` + desc(svcName) + `.`,
		RunE: release.Env1Call(func(env, timeTag string) error {
			if timeTag == `` {
				timeTag = conf.TimeTag(env)
				if err := images.Build(svcName, env, timeTag, true); err != nil {
					return err
				}
				if err := images.Push(svcName, env, timeTag); err != nil {
					return err
				}
			}
			return deploy(svcName, env, timeTag, filter)
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

func psCmdFor(svcName string) *cobra.Command {
	var filter string
	var watch bool
	cmd := &cobra.Command{
		Use:   `ps [<env>]`,
		Short: `list tasks of the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return ps(svcName, env, filter, watch)
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
			return logs(svcName, env, filter)
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
