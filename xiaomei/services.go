package main

import (
	"github.com/lovego/xiaomei/xiaomei/access"
	"github.com/lovego/xiaomei/xiaomei/deploy"
	"github.com/lovego/xiaomei/xiaomei/godoc"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/images/app"
	"github.com/lovego/xiaomei/xiaomei/images/logc"
	"github.com/lovego/xiaomei/xiaomei/images/web"
	"github.com/lovego/xiaomei/xiaomei/oam"
	"github.com/spf13/cobra"
)

func serviceCmds() []*cobra.Command {
	return append([]*cobra.Command{
		serviceCmd(`app`, `[service] the app server.`, app.Cmds()),
		serviceCmd(`web`, `[service] the web server.`, web.Cmds()),
		serviceCmd(`logc`, `[service] the log collector.`, logc.Cmds()),
		godoc.Cmd(),
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
