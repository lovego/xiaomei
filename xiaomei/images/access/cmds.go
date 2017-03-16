package access

import (
	"fmt"
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/cluster"
	"github.com/bughou-go/xiaomei/xiaomei/release"
	"github.com/bughou-go/xiaomei/xiaomei/z"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `access`,
		Short: `the access server.`,
	}
	cmd.AddCommand(setupCmd(), printCmd())
	return cmd
}

func printCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `print`,
		Short: `print access config for the project.`,
		RunE: z.NoArgCall(func() error {
			if conf, err := Config(); err != nil {
				return err
			} else {
				println(conf)
				return nil
			}
		}),
	}
}

func setupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `setup`,
		Short: `setup access config for the project.`,
		RunE: z.NoArgCall(func() error {
			if conf, err := Config(); err != nil {
				return err
			} else {
				return setupConf(conf)
			}
		}),
	}
}

func setupConf(conf string) error {
	script := fmt.Sprintf(`
	sudo tee /etc/nginx/sites-enabled/%s.conf > /dev/null &&
	sudo mkdir -p /var/log/nginx/%s &&
	sudo nginx -t &&
	sudo service nginx restart
	`, release.Name(), release.Name())
	return cluster.AccessNodeRun(cmd.O{Stdin: strings.NewReader(conf)}, script)
}
