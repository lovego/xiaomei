package cli

import (
	"errors"

	"github.com/bughou-go/xiaomei/cli/deploy"
	"github.com/bughou-go/xiaomei/cli/develop"
	"github.com/bughou-go/xiaomei/cli/setup"
	/*
		"github.com/bughou-go/xiaomei/cli/db"
		"github.com/bughou-go/xiaomei/cli/oam"
	*/

	"github.com/spf13/cobra"
)

func Run() {
	cobra.EnableCommandSorting = false

	root := cobra.Command{
		Use: `xiaomei `,
		// Short: `be small and beautiful.`,
	}
	root.AddCommand(allCmds()...)
	root.Execute()
}

func checkArgsNum(c *cobra.Command, args []string) {
	// c.Use
}

var argsCountError = errors.New(`args count doesn't match.`)

func allCmds() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   `new <project-name>`,
			Short: `create a new project.`,
			RunE: func(c *cobra.Command, args []string) error {
				switch len(args) {
				case 0:
					return errors.New(`<project-name> is required.`)
				case 1:
					return develop.New(args[0])
				default:
					return errors.New(`redundant args.`)
				}
			},
		},
		{
			Use:   `run`,
			Short: `build the binary and run it.`,
			RunE: func(c *cobra.Command, args []string) error {
				return develop.Run()
			},
		},
		{
			Use:   `build`,
			Short: `build the binary, check coding spec, compile assets.`,
			RunE: func(c *cobra.Command, args []string) error {
				return develop.Build()
			},
		},
		{
			Use:   `spec`,
			Short: `check coding specification.`,
			RunE: func(c *cobra.Command, args []string) error {
				return develop.Spec(args[0])
			},
		},
		{
			Use:   `assets`,
			Short: `compile assets.`,
			RunE: func(c *cobra.Command, args []string) error {
				develop.Assets(args)
				return nil
			},
		},
		{
			Use:   `deploy`,
			Short: `deploy project to a environment.`,
			RunE: func(c *cobra.Command, args []string) error {
				return deploy.Deploy(args[0])
			},
		},
		{
			Use:   `doc-server`,
			Short: `start doc-server.`,
			RunE: func(c *cobra.Command, args []string) error {
				return setup.DocServer()
			},
		},

		/*
			{
				Name: `restart`,
				RunE: func(c *cobra.Command, args []string) error {
					oam.Restart
				},
			},
			{
				Name: `status`,
				RunE: func(c *cobra.Command, args []string) error {
					oam.Status
				},
			},
			{
				Name: `shell`,
				RunE: func(c *cobra.Command, args []string) error {
					oam.Shell
				},
			},
			{
				Name: `exec`,
				RunE: func(c *cobra.Command, args []string) error {
					oam.Exec
				},
			},

			{
				Name: `setup`,
				RunE: func(c *cobra.Command, args []string) error {
					setup.Setup
				},
			},
			{
				Name: `setup-appserver`,
				RunE: func(c *cobra.Command, args []string) error {
					setup.SetupAppServer
				},
			},
			{
				Name: `setup-nginx`,
				RunE: func(c *cobra.Command, args []string) error {
					setup.SetupNginx
				},
			},
			{
				Name: `setup-cron`,
				RunE: func(c *cobra.Command, args []string) error {
					setup.SetupCron
				},
			},
			{
				Name: `setup-hosts`,
				RunE: func(c *cobra.Command, args []string) error {
					setup.SetupHosts
				},
			},

			{
				Name: `mysql`,
				RunE: func(c *cobra.Command, args []string) error {
					db.Mysql
				},
			},
			{
				Name: `mysqldump`,
				RunE: func(c *cobra.Command, args []string) error {
					db.Mysqldump
				},
			},
		*/
	}
}
