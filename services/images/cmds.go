package images

import (
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func Cmds(svcName string) []*cobra.Command {
	return []*cobra.Command{
		buildCmdFor(svcName),
		pushCmdFor(svcName),
		imagesCmdFor(svcName),
	}
}

func buildCmdFor(svcName string) *cobra.Command {
	var pull bool
	cmd := &cobra.Command{
		Use:   `build [<env>]`,
		Short: `[image] Build ` + imageDesc(svcName, `for`) + `.`,
		RunE: release.EnvCall(func(env string) error {
			return Build(svcName, env, release.TimeTag(env), pull)
		}),
	}
	cmd.Flags().BoolVarP(&pull, `pull`, `p`, true, `pull base image.`)
	return cmd
}

func pushCmdFor(svcName string) *cobra.Command {
	return &cobra.Command{
		Use:   `push [<env> [<tag>]]`,
		Short: `[image] Push  ` + imageDesc(svcName, `of `) + `.`,
		RunE: release.Env1Call(func(env, timeTag string) error {
			return Push(svcName, env, timeTag)
		}),
	}
}

func imagesCmdFor(svcName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `images [<env>]`,
		Short: `[image] List  images of  the ` + desc(svcName) + ` on this machine.`,
		RunE: release.EnvCall(func(env string) error {
			return List(svcName, env)
		}),
	}
	return cmd
}

func imageDesc(svcName, preposition string) string {
	if svcName == `` {
		return `images ` + preposition + ` the project`
	} else {
		return `image  ` + preposition + ` the ` + svcName + ` service`
	}
}

func desc(svcName string) string {
	if svcName == `` {
		return `project`
	} else {
		return svcName + ` service`
	}
}
