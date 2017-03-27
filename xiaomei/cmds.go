package main

import (
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/stack"
	"github.com/lovego/xiaomei/xiaomei/z"
	"github.com/spf13/cobra"
)

// Run, Build, Push, Deploy and Ps commands
func commonCmds(svcName string) []*cobra.Command {
	if svcName == `` {
		return []*cobra.Command{
			buildCmd(svcName, `all images`),
			pushCmd(svcName, `all images`),
			deployCmd(svcName, `stack`),
			psCmd(svcName, `stack`),
			logsCmd(svcName, `stack`),
		}
	} else {
		cmds := []*cobra.Command{}
		if images.Has(svcName) {
			cmds = append(cmds,
				runCmd(svcName),
				buildCmd(svcName, `the `+svcName+` image`),
				pushCmd(svcName, `the `+svcName+` image`),
			)
		}
		cmds = append(cmds,
			deployCmd(svcName, svcName+` service`),
			psCmd(svcName, svcName+` service`),
			logsCmd(svcName, svcName+` service`),
		)
		return cmds
	}
}

func runCmd(svcName string) *cobra.Command {
	var publish []string
	cmd := &cobra.Command{
		Use:   `run`,
		Short: `run    the ` + svcName + ` image.`,
		RunE: z.NoArgCall(func() error {
			return images.Run(svcName, publish)
		}),
	}
	cmd.Flags().StringSliceVarP(&publish, `publish`, `p`, nil, `publish ports for container.`)
	return cmd
}
func buildCmd(svcName, desc string) *cobra.Command {
	var pull bool
	cmd := &cobra.Command{
		Use:   `build`,
		Short: `build  ` + desc + `.`,
		RunE: z.NoArgCall(func() error {
			return images.Build(svcName, pull)
		}),
	}
	cmd.Flags().BoolVarP(&pull, `pull`, `p`, true, `pull base image.`)
	return cmd
}
func pushCmd(svcName, desc string) *cobra.Command {
	return &cobra.Command{
		Use:   `push`,
		Short: `push   ` + desc + `.`,
		RunE: z.NoArgCall(func() error {
			return images.Push(svcName)
		}),
	}
}
func deployCmd(svcName, desc string) *cobra.Command {
	var noBuild, noPush, rmCurrent bool
	cmd := &cobra.Command{
		Use:   `deploy`,
		Short: `deploy the ` + desc + `.`,
		RunE: z.NoArgCall(func() error {
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
			return stack.Deploy(svcName, rmCurrent)
		}),
	}
	cmd.Flags().BoolVarP(&noBuild, `no-build`, `B`, false, `do not build the images.`)
	cmd.Flags().BoolVarP(&noPush, `no-push`, `P`, false, `do not push the images.`)
	cmd.Flags().BoolVar(&rmCurrent, `rm-current`, false, `remove the current running `+desc+`.`)
	return cmd
}

func psCmd(svcName, desc string) *cobra.Command {
	var watch bool
	cmd := &cobra.Command{
		Use:   `ps`,
		Short: `list tasks of the ` + desc + `.`,
		RunE: func(c *cobra.Command, args []string) error {
			return stack.Ps(svcName, watch, args)
		},
	}
	cmd.Flags().BoolVarP(&watch, `watch`, `w`, false, `watch ps.`)
	return cmd
}

func logsCmd(svcName, desc string) *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   `logs`,
		Short: `list logs  of the ` + desc + `.`,
		RunE: func(c *cobra.Command, args []string) error {
			return stack.Logs(svcName, all)
		},
	}
	cmd.Flags().BoolVarP(&all, `all`, `a`, false, `list logs of all containers.`)
	return cmd
}
