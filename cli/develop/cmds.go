package develop

import (
	"errors"

	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `new <project-name>`,
			Short: `create a new project.`,
			RunE: func(c *cobra.Command, args []string) error {
				switch len(args) {
				case 0:
					return errors.New(`<project-name> is required.`)
				case 1:
					return New(args[0])
				default:
					return errors.New(`redundant args.`)
				}
			},
		},
		{
			Use:   `run`,
			Short: `build the binary and run it.`,
			RunE: func(c *cobra.Command, args []string) error {
				return Run()
			},
		},
		{
			Use:   `build`,
			Short: `build the binary, check coding spec, compile assets.`,
			RunE: func(c *cobra.Command, args []string) error {
				return Build()
			},
		},
		{
			Use:   `spec`,
			Short: `check coding specification.`,
			RunE: func(c *cobra.Command, args []string) error {
				arg := ``
				if len(args) > 0 {
					arg = args[0]
				}
				return Spec(arg)
			},
		},
		{
			Use:   `assets`,
			Short: `compile assets.`,
			RunE: func(c *cobra.Command, args []string) error {
				Assets(args)
				return nil
			},
		},
		{
			Use:   `godoc`,
			Short: `start godoc service.`,
			RunE: func(c *cobra.Command, args []string) error {
				return Godoc()
			},
		},
	}
}
