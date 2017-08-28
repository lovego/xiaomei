package main

import (
	"github.com/lovego/xiaomei/xiaomei/access"
	"github.com/lovego/xiaomei/xiaomei/deploy"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/images/app"
	"github.com/lovego/xiaomei/xiaomei/images/godoc"
	"github.com/lovego/xiaomei/xiaomei/images/logc"
	"github.com/lovego/xiaomei/xiaomei/images/tasks"
	"github.com/lovego/xiaomei/xiaomei/images/web"
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
	theCmd := &cobra.Command{
		Use:   name,
		Short: desc,
	}
	theCmd.AddCommand(images.Cmds(name)...)
	theCmd.AddCommand(deploy.Cmds(name)...)
	theCmd.AddCommand(access.Cmds(name)...)
	if len(cmds) > 0 {
		theCmd.AddCommand(cmds...)
	}
	return theCmd
}
