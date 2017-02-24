package app

import (
	"errors"

	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `new <project-name>`,
			Short: `[develop] create a new project.`,
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
			Short: `[develop] build the binary and run it.`,
			RunE: func(c *cobra.Command, args []string) error {
				return Run()
			},
		},
		{
			Use:   `build`,
			Short: `[develop] build the binary, check coding spec, compile assets.`,
			RunE: func(c *cobra.Command, args []string) error {
				return Build()
			},
		},
		{
			Use:   `spec`,
			Short: `[develop] check coding specification.`,
			RunE: func(c *cobra.Command, args []string) error {
				arg := ``
				if len(args) > 0 {
					arg = args[0]
				}
				return Spec(arg)
			},
		},
		{
			Use:   `deps`,
			Short: `[develop] list all dependences of project.`,
			Run: func(c *cobra.Command, args []string) {
				Dependences()
			},
		},
		copy2vendorCmd(),
	}
}

func copy2vendorCmd() *cobra.Command {
	var n bool
	cmd := &cobra.Command{
		Use:   `copy2vendor`,
		Short: `[develop] copy the specified packages to project vendor dir.`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) <= 0 {
				return errors.New(`need at least a package path`)
			}
			return Copy2Vendor(args, n)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&n, `no-clobber`, `n`, false, `do not overwrite an existing file.`)
	return cmd
}
