package cluster

import (
	"fmt"
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/release"
	"github.com/bughou-go/xiaomei/xiaomei/z"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `cluster`,
		Short: `cluster operations.`,
	}
	cmd.AddCommand(lsCmd(), shellCmd())
	return cmd
}

func lsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `ls`,
		Short: `list all clusters.`,
		RunE: z.NoArgCall(func() error {
			return lsClusters()
		}),
	}
}

func shellCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `shell`,
		Short: `ssh into a manager of cluster.`,
		RunE: z.NoArgCall(func() error {
			return Run(cmd.O{}, ``)
		}),
	}
}

func lsClusters() error {
	managers, workers := release.GetCluster().List()
	cmd.Run(cmd.O{}, `echo`, fmt.Sprintf("managers: %s", strings.Join(managers, "\t")))
	cmd.Run(cmd.O{}, `echo`, fmt.Sprintf("workers: %s", strings.Join(workers, "\t")))
	return nil
}

func Run(o cmd.O, script string) error {
	_, err := cmd.SshRun(o, release.GetCluster().SshAddr(), script)
	return err
}

func SshRun(o cmd.O, script string) (string, error) {
	return cmd.SshRun(o, release.GetCluster().SshAddr(), script)
}
