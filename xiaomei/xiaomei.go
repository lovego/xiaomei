package main

import (
	"os"

	"github.com/lovego/xiaomei/xiaomei/cluster"
	// "github.com/lovego/xiaomei/xiaomei/images/access"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/images/app"
	"github.com/lovego/xiaomei/xiaomei/images/logc"
	"github.com/lovego/xiaomei/xiaomei/images/web"
	"github.com/lovego/xiaomei/xiaomei/new"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false

	appCmd := app.Cmd()
	webCmd := web.Cmd()
	// accessCmd := access.Cmd()
	logcCmd := logc.Cmd()

	appCmd.AddCommand(commonCmds(`app`)...)
	webCmd.AddCommand(commonCmds(`web`)...)
	// accessCmd.AddCommand(commonCmds(`access`)...)
	logcCmd.AddCommand(commonCmds(`logc`)...)

	root := rootCmd()
	root.AddCommand(appCmd, webCmd /*accessCmd,*/, logcCmd, cluster.Cmd())
	root.AddCommand(
		buildCmdFor(``),
		pushCmdFor(``),
		deployCmdFor(``),
		psCmdFor(``),
		logsCmdFor(``),
	)
	root.AddCommand(new.Cmd(), yamlCmd(), versionCmd())
	root.Execute()
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `xiaomei`,
		Short: `be small and beautiful.`,
	}

	if release.Arg1IsEnv() {
		cmd.SetArgs(os.Args[2:])
	}
	return cmd
}

func commonCmds(svcName string) (cmds []*cobra.Command) {
	if images.Has(svcName) {
		cmds = append(cmds,
			buildCmdFor(svcName),
			pushCmdFor(svcName),
			runCmdFor(svcName),
		)
	}
	cmds = append(cmds,
		deployCmdFor(svcName),
		psCmdFor(svcName),
		logsCmdFor(svcName),
	)
	return
}
