package main

import (
	"fmt"

	"github.com/bughou-go/xiaomei/xiaomei/images"
	"github.com/bughou-go/xiaomei/xiaomei/stack"
	"github.com/bughou-go/xiaomei/xiaomei/z"
	"github.com/spf13/cobra"
)

// Run, Build, Push, Deploy and Ps commands
func commonCmds(svcName string) []*cobra.Command {
	cmds := []*cobra.Command{}
	if svcName != `` {
		cmds = append(cmds, runCmd(svcName))
	}
	var target, s string
	if svcName == `` {
		target, s = `all`, `s`
	} else {
		target, s = svcName, ``
	}
	cmds = append(cmds,
		buildCmd(svcName, target, s),
		pushCmd(svcName, target, s),
		deployCmd(svcName, target, s),
		psCmd(svcName, target, s),
	)
	return cmds
}

func runCmd(svcName string) *cobra.Command {
	var ports []string
	cmd := &cobra.Command{
		Use:   `run`,
		Short: fmt.Sprintf(`run    %s image.`, svcName),
		RunE: z.NoArgCall(func() error {
			return images.Run(svcName /*, ports*/)
		}),
	}
	cmd.Flags().StringSliceVarP(&ports, `ports`, `p`, nil, `specify ports`)
	return cmd
}

func buildCmd(svcName, target, s string) *cobra.Command {
	return &cobra.Command{
		Use:   `build`,
		Short: fmt.Sprintf(`build  %s image%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return images.Build(svcName)
		}),
	}
}

func pushCmd(svcName, target, s string) *cobra.Command {
	return &cobra.Command{
		Use:   `push`,
		Short: fmt.Sprintf(`push   %s image%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return images.Push(svcName)
		}),
	}
}

func deployCmd(svcName, target, s string) *cobra.Command {
	var build, push bool
	cmd := &cobra.Command{
		Use:   `deploy`,
		Short: fmt.Sprintf(`deploy %s service%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return stack.Deploy(svcName, build, push)
		}),
	}
	cmd.Flags().BoolVarP(&build, `build`, `b`, false, fmt.Sprintf(`build the image%s.`, s))
	cmd.Flags().BoolVarP(&push, `push`, `p`, false, fmt.Sprintf(`push the image%s.`, s))
	return cmd
}

func psCmd(svcName, target, s string) *cobra.Command {
	return &cobra.Command{
		Use:   `ps`,
		Short: fmt.Sprintf(`list tasks of %s service%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return stack.Ps(svcName)
		}),
	}
}
