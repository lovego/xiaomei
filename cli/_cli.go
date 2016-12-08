package cli

import (
	"os"

	"github.com/bughou-go/xiaomei/cli/develop"
	/*
		"github.com/bughou-go/xiaomei/cli/db"
		"github.com/bughou-go/xiaomei/cli/deploy"
		"github.com/bughou-go/xiaomei/cli/oam"
		"github.com/bughou-go/xiaomei/cli/setup"
	*/

	"github.com/urfave/cli"
)

func Run() {
	app := cli.NewApp()
	app.Name = `xiaomei`
	app.Usage = `be small and beautiful.`
	app.Version = `0.0.1`
	app.Commands = allCmds()
	app.Run(os.Args)
}

func allCmds() []cli.Command {
	return []cli.Command{
		{
			Name:      `new`,
			Usage:     `create a new project.`,
			ArgsUsage: `<project-name>`,
			Category:  `develop`,
			Action: func(c *cli.Context) error {
				return develop.New(c.Args().First())
			},
		},
		{
			Name:     `run`,
			Usage:    `build the binary and run it.`,
			Category: `develop`,
			Action: func(c *cli.Context) error {
				return develop.Run()
			},
		},
		{
			Name:     `build`,
			Usage:    `build the binary, check code spec, compile assets.`,
			Category: `develop`,
			Action: func(c *cli.Context) error {
				return develop.Build()
			},
		},
		{
			Name:     `spec`,
			Usage:    `build a project.`,
			Category: `develop`,
			Action: func(c *cli.Context) error {
				return develop.Spec(c.Args().First())
			},
		},
		{
			Name:     `assets`,
			Category: `develop`,
			Action: func(c *cli.Context) error {
				develop.Assets(c.Args())
				return nil
			},
		},
		/*
			{
				Name: `deploy`,
				Action: func(c *cli.Context) error {
					deploy.Deploy
				},
			},
			{
				Name: `update-doc`,
				Action: func(c *cli.Context) error {
					deploy.UpdateDoc
				},
			},

			{
				Name: `restart`,
				Action: func(c *cli.Context) error {
					oam.Restart
				},
			},
			{
				Name: `status`,
				Action: func(c *cli.Context) error {
					oam.Status
				},
			},
			{
				Name: `shell`,
				Action: func(c *cli.Context) error {
					oam.Shell
				},
			},
			{
				Name: `exec`,
				Action: func(c *cli.Context) error {
					oam.Exec
				},
			},

			{
				Name: `setup`,
				Action: func(c *cli.Context) error {
					setup.Setup
				},
			},
			{
				Name: `setup-appserver`,
				Action: func(c *cli.Context) error {
					setup.SetupAppServer
				},
			},
			{
				Name: `setup-nginx`,
				Action: func(c *cli.Context) error {
					setup.SetupNginx
				},
			},
			{
				Name: `setup-cron`,
				Action: func(c *cli.Context) error {
					setup.SetupCron
				},
			},
			{
				Name: `setup-hosts`,
				Action: func(c *cli.Context) error {
					setup.SetupHosts
				},
			},

			{
				Name: `mysql`,
				Action: func(c *cli.Context) error {
					db.Mysql
				},
			},
			{
				Name: `mysqldump`,
				Action: func(c *cli.Context) error {
					db.Mysqldump
				},
			},
		*/
	}
}
