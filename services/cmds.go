package services

import (
	"github.com/lovego/xiaomei/services/deploy"
	"github.com/lovego/xiaomei/services/images"
	"github.com/lovego/xiaomei/services/images/app"
	"github.com/lovego/xiaomei/services/images/logc"
	"github.com/lovego/xiaomei/services/images/web"
	"github.com/lovego/xiaomei/services/oam"
	"github.com/lovego/xiaomei/services/run"
	"github.com/spf13/cobra"
)

func Cmds() []*cobra.Command {
	cmds := []*cobra.Command{
		serviceCmd(`app`, `[service] the app server.`, app.Cmds()),
		serviceCmd(`web`, `[service] the web server.`, web.Cmds()),
		serviceCmd(`logc`, `[service] the log collector.`, logc.Cmds()),
	}
	cmds = append(cmds, deploy.Cmds(``)...)
	cmds = append(cmds, oam.Cmds(``)...)
	cmds = append(cmds, images.Cmds(``)...)
	return cmds
}

func serviceCmd(name, desc string, cmds []*cobra.Command) *cobra.Command {
	theCmd := &cobra.Command{
		Use:   name,
		Short: desc,
	}
	theCmd.AddCommand(run.Cmds(name)...)
	theCmd.AddCommand(deploy.Cmds(name)...)
	theCmd.AddCommand(oam.Cmds(name)...)
	theCmd.AddCommand(images.Cmds(name)...)
	if len(cmds) > 0 {
		theCmd.AddCommand(cmds...)
	}
	return theCmd
}
