package stack

import (
	"fmt"

	"github.com/bughou-go/xiaomei/cli/z"
	"github.com/spf13/cobra"
)

// Run, Build, Push, Deploy and Ps commands
func Cmds(svcName string) []*cobra.Command {
	var target, s string
	if svcName == `` {
		target, s = `all`, `s`
	} else {
		target, s = svcName, ``
	}
	cmds := []*cobra.Command{
		buildCmd(svcName, target, s),
		pushCmd(svcName, target, s),
		deployCmd(svcName, target, s),
		psCmd(svcName, target, s),
	}
	if svcName != `` {
		cmds = append(cmds, runCmd(svcName))
	}
	return cmds
}

func runCmd(svcName string) *cobra.Command {
	return &cobra.Command{
		Use:   `run`,
		Short: fmt.Sprintf(`run %s image.`, svcName),
		RunE: z.NoArgCall(func() error {
			return Build(svcName)
		}),
	}
}

func buildCmd(svcName, target, s string) *cobra.Command {
	return &cobra.Command{
		Use:   `build`,
		Short: fmt.Sprintf(`build %s image%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return Build(svcName)
		}),
	}
}

func pushCmd(svcName, target, s string) *cobra.Command {
	return &cobra.Command{
		Use:   `push`,
		Short: fmt.Sprintf(`push %s image%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return Push(svcName)
		}),
	}
}

func deployCmd(svcName, target, s string) *cobra.Command {
	var doBuild, doPush *bool
	cmd := &cobra.Command{
		Use:   `deploy`,
		Short: fmt.Sprintf(`deploy %s service%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return Deploy(svcName, *doBuild, *doPush)
		}),
	}
	doBuild = cmd.Flags().BoolP(`build`, `b`, false, fmt.Sprintf(`build the image%s.`, s))
	doPush = cmd.Flags().BoolP(`push`, `p`, false, fmt.Sprintf(`push the image%s.`, s))
	return cmd
}

func psCmd(svcName, target, s string) *cobra.Command {
	return &cobra.Command{
		Use:   `ps`,
		Short: fmt.Sprintf(`list tasks of %s service%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return Ps(svcName)
		}),
	}
}
