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
			RunE: func(c *cobra.Command, args []string) error {
				return build(svcName)
			},
		},
		{
			Use:   `deploy <env>`,
			Short: fmt.Sprintf(`deploy %s service%s.`, target, s),
			RunE: z.Arg1Call(`dev`, func(env string) error {
				return deploy(env, svcName)
			}),
		},
		{
			Use:   `ps [<env>]`,
			Short: fmt.Sprintf(`list tasks of %s service%s.`, target, s),
			RunE: z.Arg1Call(`dev`, func(env string) error {
				return ps(env, svcName)
			}),
		},
	}
}
