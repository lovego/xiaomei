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
		Use:   `run`,
		Short: `run    the ` + desc(svcName) + `.`,
		RunE: release.NoArgCall(func() error {
			return run(svcName)
		}),
	}
	// cmd.Flags().StringSliceVarP(&publish, `publish`, `p`, nil, `publish ports for container.`)
	return cmd
}

func deployCmdFor(svcName string) *cobra.Command {
	var noBuild, noPush /*, rmCurrent*/ bool
	var filter string
	cmd := &cobra.Command{
		Use:   `deploy`,
		Short: `deploy the ` + desc(svcName) + `.`,
		RunE: release.NoArgCall(func() error {
			if !noBuild {
				if err := images.Build(svcName, true); err != nil {
					return err
				}
			}
			if !noPush {
				if err := images.Push(svcName); err != nil {
					return err
				}
			}
			return getDriver().Deploy(svcName, filter) // , rmCurrent)
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
		Use:   `rm-deploy`,
		Short: `remove deployment of the ` + desc(svcName) + `.`,
		RunE: release.NoArgCall(func() error {
			return getDriver().RmDeploy(svcName, filter)
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	return cmd
}

func psCmdFor(svcName string) *cobra.Command {
	var watch bool
	var filter string
	cmd := &cobra.Command{
		Use:   `ps`,
		Short: `list tasks of the ` + desc(svcName) + `.`,
		RunE: func(c *cobra.Command, args []string) error {
			return getDriver().Ps(svcName, filter, watch, args)
		},
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	cmd.Flags().BoolVarP(&watch, `watch`, `w`, false, `watch ps.`)
	return cmd
}

func logsCmdFor(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:   `logs`,
		Short: `list logs  of the ` + desc(svcName) + `.`,
		RunE: func(c *cobra.Command, args []string) error {
			return getDriver().Logs(svcName, filter)
		},
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
