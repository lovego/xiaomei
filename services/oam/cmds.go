package oam

import (
	"strings"

	cmdPkg "github.com/lovego/cmd"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

// Shell, Ps, Logs, ... commands
func Cmds(svcName string) (cmds []*cobra.Command) {
	if svcName == `` {
		cmds = append(cmds, lsCmd(), shellCmd())
	} else {
		cmds = append(cmds, shellCmdFor(svcName))
	}
	cmds = append(cmds,
		psCmdFor(svcName),
		logsCmdFor(svcName),
	)
	cmds = append(cmds, operationCmdsFor(svcName)...)
	return
}

func shellCmd() *cobra.Command {
	var filter string
	theCmd := &cobra.Command{
		Use:   `shell [env]`,
		Short: `[oam] Enter a node shell.`,
		RunE: release.EnvCall(func(env string) error {
			_, err := release.GetCluster(env).Run(filter, cmdPkg.O{}, ``)
			return err
		}),
	}
	theCmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr`)
	return theCmd
}

func lsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `ls`,
		Short: `[oam] List all nodes.`,
		RunE: release.EnvCall(func(env string) error {
			release.GetCluster(env).List()
			return nil
		}),
	}
}

func shellCmdFor(svcName string) *cobra.Command {
	var filter string
	theCmd := &cobra.Command{
		Use:   `shell [env]`,
		Short: `[oam] Enter a container for ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return shell(svcName, env, filter)
		}),
	}
	theCmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr`)
	return theCmd
}

func psCmdFor(svcName string) *cobra.Command {
	var filter string
	var watch bool
	cmd := &cobra.Command{
		Use:   `ps [env]`,
		Short: `[oam] List containers of the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return ps(svcName, env, filter, watch)
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	cmd.Flags().BoolVarP(&watch, `watch`, `w`, false, `watch ps.`)
	return cmd
}

func logsCmdFor(svcName string) *cobra.Command {
	var filter string
	cmd := &cobra.Command{
		Use:                   `logs [env] [-- "docker logs" flags]`,
		DisableFlagsInUseLine: true,
		Short:                 `[oam] List logs  of the ` + desc(svcName) + `.`,
		RunE: func(c *cobra.Command, args []string) error {
			var env, options string
			if len(args) > 0 {
				env = args[0]
			}
			env, err := release.CheckEnv(env)
			if err != nil {
				return err
			}
			if len(args) > 1 {
				options = strings.Join(args[1:], " ")
			}

			return logs(svcName, env, filter, options)
		},
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	return cmd
}

func operationCmdsFor(svcName string) []*cobra.Command {
	var operations = []string{`start`, `stop`, `restart`}
	cmds := make([]*cobra.Command, len(operations), len(operations))
	for i, operation := range operations {
		cmds[i] = makeOperationCmd(operation, svcName)
	}
	return cmds
}

func makeOperationCmd(operation, svcName string) *cobra.Command {
	var filter, startTimeout string
	cmd := &cobra.Command{
		Use:   operation + ` [env]`,
		Short: `[oam] ` + strings.Title(operation) + ` the ` + desc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return operate(operation, svcName, env, filter, startTimeout)
		}),
	}
	cmd.Flags().StringVarP(&filter, `filter`, `f`, ``, `filter by node addr.`)
	switch operation {
	case "start", "restart":
		cmd.Flags().StringVarP(&startTimeout, `start-timeout`, `t`, `1m`,
			`Timeout to wait for starting (until a " started." is printed on stdout or stderr).`)
	}

	return cmd
}

func desc(svcName string) string {
	if svcName == `` {
		return `project`
	} else {
		return svcName + ` service`
	}
}
