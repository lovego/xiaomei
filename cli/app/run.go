package app

import (
	"github.com/bughou-go/xiaomei/cli/app/deps"
	"github.com/bughou-go/xiaomei/cli/stack"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `app`,
		Short: `the appserver.`,
	}
	cmd.AddCommand(
		RunCmd(),
		BuildCmd(),
		CheckCodeCmd(),
		deps.Cmd(),
	)
	return cmd
}

func RunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `run`,
		Short: `build the binary and run it.`,
		RunE: func(c *cobra.Command, args []string) error {
			return Run()
		},
	}
}

func Run() error {
	if err := BuildBinary(); err != nil {
		return err
	}

	var image string
	if svc, err := stack.GetService(`app`); err != nil {
		return err
	} else if image, err = svc.GetImage(); err != nil {
		return err
	}

	_, err := cmd.Run(cmd.O{}, `docker`,
		`run`, `--name=`+config.DeployName(), `-it`, `--rm`, `--network=host`,
		`-v`, config.Root()+`:/home/ubuntu/appserver`,
		image,
	)
	return err
}
