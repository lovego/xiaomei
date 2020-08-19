package services

import (
	"github.com/lovego/xiaomei/services/deploy"
	"github.com/lovego/xiaomei/services/images"
	"github.com/lovego/xiaomei/services/images/app"
	"github.com/lovego/xiaomei/services/oam"
	"github.com/lovego/xiaomei/services/run"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	cmds := []*cobra.Command{
		serviceCmd(`app`, `[service] The app server.`),
		serviceCmd(`web`, `[service] The web server.`),
		serviceCmd(`logc`, `[service] The log collector.`),
	}
	cmds = append(cmds, app.Cmds()...)
	cmds = append(cmds, deploy.Cmds(``)...)
	cmds = append(cmds, oam.Cmds(``)...)
	cmds = append(cmds, images.Cmds(``)...)
	return cmds
}

func serviceCmd(name, desc string) *cobra.Command {
	theCmd := &cobra.Command{
		Use:   name,
		Short: desc,
	}
	theCmd.AddCommand(run.Cmds(name)...)
	theCmd.AddCommand(deploy.Cmds(name)...)
	theCmd.AddCommand(oam.Cmds(name)...)
	theCmd.AddCommand(images.Cmds(name)...)
	return theCmd
}
