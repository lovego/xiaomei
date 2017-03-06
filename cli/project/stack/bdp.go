package stack

import (
	"fmt"

	"github.com/bughou-go/xiaomei/cli/z"
	"github.com/spf13/cobra"
)

// Build, Deploy and Ps commands
func BDPcmds(svcName string) []*cobra.Command {
	var target, s string
	if svcName == `` {
		target, s = `all`, `s`
	} else {
		target, s = svcName, ``
	}

	return []*cobra.Command{
		{
			Use:   `build`,
			Short: fmt.Sprintf(`build %s image%s.`, target, s),
			RunE: z.NoArgCall(func() error {
				return Build(svcName)
			}),
		},
		deployCmd(svcName, target, s),
		{
			Use:   `ps`,
			Short: fmt.Sprintf(`list tasks of %s service%s.`, target, s),
			RunE: z.NoArgCall(func() error {
				return Ps(svcName)
			}),
		},
	}
}

func deployCmd(svcName, target, s string) *cobra.Command {
	var noBuild, noPush *bool
	cmd := &cobra.Command{
		Use:   `deploy`,
		Short: fmt.Sprintf(`deploy %s service%s.`, target, s),
		RunE: z.NoArgCall(func() error {
			return Deploy(svcName, *noBuild, *noPush)
		}),
	}
	noBuild = cmd.Flags().BoolP(`no-build`, `B`, false, fmt.Sprintf(`do not build the image%s.`, s))
	noPush = cmd.Flags().BoolP(`no-push`, `P`, false, fmt.Sprintf(`do not build the image%s.`, s))
	return cmd
}
