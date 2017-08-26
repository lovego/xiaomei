package main

import (
	"errors"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/access"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/images/app"
	"github.com/lovego/xiaomei/xiaomei/images/godoc"
	"github.com/lovego/xiaomei/xiaomei/images/logc"
	"github.com/lovego/xiaomei/xiaomei/images/tasks"
	"github.com/lovego/xiaomei/xiaomei/images/web"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func serviceCmds() []*cobra.Command {
	return []*cobra.Command{
		serviceCmd(`app`, `the app server.`, app.Cmds()),
		serviceCmd(`tasks`, `the app tasks.`, tasks.Cmds()),
		serviceCmd(`web`, `the web server.`, web.Cmds()),
		serviceCmd(`logc`, `the log collector.`, logc.Cmds()),
		serviceCmd(`godoc`, `the godoc server.`, godoc.Cmds()),
	}
}

func serviceCmd(name, desc string, cmds []*cobra.Command) *cobra.Command {
	var filter string
	theCmd := &cobra.Command{
		Use:   name,
		Short: desc,
		RunE: func(thisCmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New(`redundant args.`)
			}
			if release.Arg1IsEnv() {
				_, err := cluster.ServiceRun(cmd.O{}, name, filter,
					`docker exec -it --detach-keys='ctrl-@' `+conf.ContainerNameOf(name)+` bash`,
				)
				return err
			} else {
				return thisCmd.Help()
			}
		},
	}
	theCmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr`)
	if len(cmds) > 0 {
		theCmd.AddCommand(cmds...)
	}
	theCmd.AddCommand(images.Cmds(name)...)
	theCmd.AddCommand(deploy.Cmds(name)...)
	theCmd.AddCommand(access.Cmds(name)...)
	return theCmd
}
