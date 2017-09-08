package images

import (
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
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
	var tag, pull bool
	cmd := &cobra.Command{
		Use:   `build [<env>]`,
		Short: `build  ` + imageDesc(svcName) + `.`,
		RunE: release.EnvCall(func(env string) error {
			timeTag := ``
			if tag {
				timeTag = conf.TimeTag(env)
			}
			return Build(svcName, env, timeTag, pull)
		}),
	}
	cmd.Flags().BoolVarP(&tag, `tag`, `t`, false, `add a deploy time tag.`)
	cmd.Flags().BoolVarP(&pull, `pull`, `p`, true, `pull base image.`)
	return cmd
}

func pushCmdFor(svcName string) *cobra.Command {
	return &cobra.Command{
		Use:   `push [<env> [<tag>]]`,
		Short: `push   ` + imageDesc(svcName) + `.`,
		RunE: release.Env1Call(func(env, timeTag string) error {
			return Push(svcName, env, timeTag)
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
