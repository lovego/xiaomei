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
	"github.com/lovego/xiaomei/xiaomei/oam"
	"github.com/lovego/xiaomei/xiaomei/workspace_godoc"
	"github.com/spf13/cobra"
)

func serviceCmds() []*cobra.Command {
	return append([]*cobra.Command{
		serviceCmd(`app`, `[service] the app server.`, app.Cmds()),
		serviceCmd(`tasks`, `[service] the app tasks.`, tasks.Cmds()),
		serviceCmd(`web`, `[service] the web server.`, web.Cmds()),
		serviceCmd(`logc`, `[service] the log collector.`, logc.Cmds()),
		serviceCmd(`godoc`, `[service] the godoc server.`, godoc.Cmds()),
		workspace_godoc.Cmd(),
	})
}

func serviceCmd(name, desc string, cmds []*cobra.Command) *cobra.Command {
	theCmd := &cobra.Command{
		Use:   name,
		Short: desc,
	}
	theCmd.AddCommand(deploy.Cmds(name)...)
	if cmd := access.Cmd(name); cmd != nil {
		theCmd.AddCommand(cmd)
	}
	theCmd.AddCommand(oam.Cmds(name)...)
	theCmd.AddCommand(images.Cmds(name)...)
	if len(cmds) > 0 {
		theCmd.AddCommand(cmds...)
	}
	return theCmd
}
