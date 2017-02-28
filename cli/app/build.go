package app

import (
	"errors"
	"path/filepath"

	"github.com/bughou-go/xiaomei/cli/project"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/spf13/cobra"
)

func BuildCmd() *cobra.Command {
	var binary, spec, assets, image bool
	cmd := &cobra.Command{
		Use:   `build`,
		Short: `build binary, check code spec, build assets, build docker image.`,
		RunE: func(c *cobra.Command, args []string) error {
			return Build(binary, spec, assets, image)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&binary, `binary`, `b`, false, `build binary.`)
	flags.BoolVarP(&spec, `spec`, `s`, false, `check code spec.`)
	flags.BoolVarP(&assets, `assets`, `a`, false, `build assets.`)
	flags.BoolVarP(&image, `image`, `i`, false, `build image.`)
	return cmd
}

func Build(binary, spec, assets, image bool) error {
	if !(binary || spec || assets || image) {
		binary, spec, assets, image = true, true, true, true
	}
	if binary {
		if err := BuildBinary(); err != nil {
			return err
		}
	}
	if spec {
		if err := Spec(``); err != nil {
			return err
		}
	}
	if assets {
		/*
			if err :=	Assets(nil); err != nil {
				return err
			}
		*/
	}
	if image {
		return BuildImage()
	}
	return nil
}

func BuildBinary() error {
	config.Log(`building binary.`)
	if cmd.Ok(cmd.O{
		Dir: filepath.Join(config.Root(), "../.."),
		Env: []string{`GOBIN=` + config.Root()},
	}, `go`, `install`) {
		return nil
	}
	return errors.New(`building binary failed.`)
}

func BuildImage() error {
	if svc, err := project.GetService(`app`); err != nil {
		return err
	} else if err = svc.BuildImage(); err != nil {
		return err
	}
	return nil
}
