package cluster

import (
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/slice"
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
			managers, workers := GetCluster().List()
			println(`managers: `, strings.Join(managers, "\t"))
			println(`workers: `, strings.Join(workers, "\t"))
			return nil
		}),
	}
}

func shellCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `shell`,
		Short: `ssh into a manager of cluster.`,
		RunE: z.NoArgCall(func() error {
			_, err := Run(cmd.O{}, ``)
			return err
		}),
	}
}

func Run(o cmd.O, script string) (string, error) {
	return cmd.SshRun(o, GetCluster().SshAddr(), script)
}

func AccessNodeRun(o cmd.O, script string) error {
	cluster := GetCluster()
	if err := accessNodeRun(o, script, cluster.Managers); err != nil {
		return err
	}
	return accessNodeRun(o, script, cluster.Workers)
}

func accessNodeRun(o cmd.O, script string, nodes []Node) error {
	for _, node := range nodes {
		if slice.ContainsString(node.Labels, `hasAccess=true`) {
			if _, err := cmd.SshRun(o, node.SshAddr(), script); err != nil {
				return err
			}
		}
	}
	return nil
}
