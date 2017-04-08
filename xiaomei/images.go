package main

import (
	"github.com/lovego/xiaomei/xiaomei/deploy"
	"github.com/lovego/xiaomei/xiaomei/images"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

// Build, Push commands
func buildCmdFor(svcName string) *cobra.Command {
	var pull bool
	cmd := &cobra.Command{
		Use:   `build`,
		Short: `build  ` + imageDesc(svcName) + `.`,
		RunE: release.NoArgCall(func() error {
			return buildImage(svcName, pull)
		}),
	}
	cmd.Flags().BoolVarP(&pull, `pull`, `p`, true, `pull base image.`)
	return cmd
}

func pushCmdFor(svcName string) *cobra.Command {
	return &cobra.Command{
		Use:   `push`,
		Short: `push   ` + imageDesc(svcName) + `.`,
		RunE: release.NoArgCall(func() error {
			return pushImage(svcName)
		}),
	}
}

func buildImage(svcName string, pull bool) error {
	if svcName == `` {
		return eachServiceDo(func(svcName string) error {
			return images.Build(svcName, deploy.ImageNameOf(svcName), pull)
		})
	}
	return images.Build(svcName, deploy.ImageNameOf(svcName), pull)
}

func pushImage(svcName string) error {
	if svcName == `` {
		return eachServiceDo(func(svcName string) error {
			return images.Push(svcName, deploy.ImageNameOf(svcName))
		})
	}
	return images.Push(svcName, deploy.ImageNameOf(svcName))
}

func imageDesc(svcName string) string {
	if svcName == `` {
		return `all images`
	} else {
		return `the ` + svcName + ` image`
	}
}

func eachServiceDo(work func(svcName string) error) error {
	for svcName := range deploy.ServiceNames() {
		if svcName != `` {
			if err := work(svcName); err != nil {
				return err
			}
		}
	}
	return nil
}
