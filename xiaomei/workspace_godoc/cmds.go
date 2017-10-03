package workspace_godoc

import (
	"fmt"
	"os"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `workspace-godoc`,
		Short: `deploy the workspace godoc server.`,
		RunE:  release.NoArgCall(deploy),
	}
	cmd.AddCommand(shellCmd())
	return cmd
}

func deploy() error {
	script := fmt.Sprintf(`
  docker stop workspace-godoc >/dev/null 2>&1 && docker rm workspace-godoc
	docker run --name=workspace-godoc -d --restart=always \
	--network=host -e=GODOCPORT=1234 \
	-v %s:/home/ubuntu/go hub.c.163.com/lovego/xiaomei/godoc
	`, os.Getenv(`GOPATH`))
	_, err := cmd.Run(cmd.O{}, `sh`, `-c`, script)
	return err
}

func shellCmd() *cobra.Command {
	theCmd := &cobra.Command{
		Use:   `shell [<env>]`,
		Short: `enter the container for workspace godoc server.`,
		RunE:  release.NoArgCall(shell),
	}
	return theCmd
}

func shell() error {
	_, err := cmd.Run(cmd.O{}, `docker`, `exec`, `-it`,
		`--detach-keys=ctrl-@`, `workspace-godoc`, `bash`,
	)
	return err
}
