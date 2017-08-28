package main

import (
	"fmt"

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
	theCmd := &cobra.Command{
		Use:   name,
		Short: desc,
	}
	theCmd.AddCommand(images.Cmds(name)...)
	theCmd.AddCommand(deploy.Cmds(name)...)
	theCmd.AddCommand(access.Cmds(name)...)
	theCmd.AddCommand(serviceShellCmd(name))
	if len(cmds) > 0 {
		theCmd.AddCommand(cmds...)
	}
	return theCmd
}

func serviceShellCmd(svcName string) *cobra.Command {
	var filter string
	theCmd := &cobra.Command{
		Use:   `shell [<env>]`,
		Short: fmt.Sprintf(`enter a container for %s`, svcName),
		RunE: release.EnvCall(func(env string) error {
			_, err := cluster.Get(env).ServiceRun(svcName, filter, cmd.O{},
				`docker exec -it --detach-keys='ctrl-@' `+conf.FirstContainerNameOf(env, svcName)+` bash`,
			)
			return err
		}),
	}
	theCmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr`)
	return theCmd
}
