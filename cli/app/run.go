package app

import (
	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/spf13/cobra"
)

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

	image, err := stack.ImageName(`app`)
	if err != nil {
		return err
	}

	_, err = cmd.Run(cmd.O{}, `docker`,
		`run`, `--name=`+config.DeployName(), `-it`, `--rm`, `--network=host`,
		`-v`, config.Root()+`:/home/ubuntu/appserver`, image,
	)
	return err
}
