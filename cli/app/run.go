package app

import (
	"github.com/bughou-go/xiaomei/cli/project/stack"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `run`,
		Short: `build the binary and run it.`,
		RunE: func(c *cobra.Command, args []string) error {
			return run()
		},
	}
}

func run() error {
	imageName, err := stack.ImageName(`app`)
	if err != nil {
		return err
	}

	if cmd.Ok(cmd.O{NoStdout: true, NoStderr: true}, `docker`, `image`, `inspect`, imageName) {
		if err := buildBinary(); err != nil {
			return err
		}
	} else {
		if err := stack.Build(`app`); err != nil {
			return err
		}
	}

	_, err = cmd.Run(cmd.O{}, `docker`,
		`run`, `--name=`+config.DeployName(), `-it`, `--rm`, `--network=host`,
		`-v`, config.Root()+`:/home/ubuntu/appserver`, imageName,
	)
	return err
}
