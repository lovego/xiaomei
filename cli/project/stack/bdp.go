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
		{
			Use:   `deploy`,
			Short: fmt.Sprintf(`deploy %s service%s.`, target, s),
			RunE: z.NoArgCall(func() error {
				return Deploy(svcName)
			}),
		},
		{
			Use:   `ps`,
			Short: fmt.Sprintf(`list tasks of %s service%s.`, target, s),
			RunE: z.NoArgCall(func() error {
				return Ps(svcName)
			}),
		},
	}
}
