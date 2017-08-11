package images

import (
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds(svcName string) []*cobra.Command {
	return []*cobra.Command{
		buildCmdFor(svcName),
		pushCmdFor(svcName),
	}
}

func buildCmdFor(svcName string) *cobra.Command {
	var pull bool
	cmd := &cobra.Command{
		Use:   `build`,
		Short: `build  ` + imageDesc(svcName) + `.`,
		RunE: release.NoArgCall(func() error {
			return Build(svcName, pull)
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
			return Push(svcName)
		}),
	}
}

func imageDesc(svcName string) string {
	if svcName == `` {
		return `all images`
	} else {
		return `the ` + svcName + ` image`
	}
}
