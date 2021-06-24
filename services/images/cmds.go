package images

import (
	"errors"

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
	var doPush bool
	var push Push
	cmd := &cobra.Command{
		Use:   `build [flags] [env] [ -- [prepare flags] [-- docker build flags] ]`,
		Short: `[image] Build ` + imageDesc(svcName, `for`) + `.`,
		RunE: release.EnvSlicesCall(func(env string, args [][]string) error {
			if len(args) > 3 || len(args) > 0 && len(args[0]) > 0 {
				return errors.New("invalid arguments usage.")
			}
			b := Build{Env: env, Tag: release.TimeTag(env)}
			if len(args) > 1 {
				b.PrepareFlags = args[1]
			}
			if len(args) > 2 {
				b.DockerBuildFlags = args[2]
			}
			if err := b.Run(svcName); err != nil {
				return err
			}
			if doPush {
				push.Env, push.Tag = b.Env, b.Tag
				return push.Run(svcName)
			}
			return nil
		}),
		DisableFlagsInUseLine: true,
	}
	cmd.Flags().BoolVarP(&doPush, "push", "P", false, "push "+imageDesc(svcName, `for`))
	cmd.Flags().StringVarP(&push.DockerLogin.User, "docker-user", "u", "",
		"The docker login user, used before pushing image")
	cmd.Flags().StringVarP(&push.DockerLogin.Password, "docker-password", "p", "",
		"The docker login password, used before pushing image")
	return cmd
}

func pushCmdFor(svcName string) *cobra.Command {
	var push Push
	cmd := &cobra.Command{
		Use:   `push [env [tag]]`,
		Short: `[image] Push  ` + imageDesc(svcName, `of `) + `.`,
		RunE: release.Env1Call(func(env, timeTag string) error {
			push.Env, push.Tag = env, timeTag
			return push.Run(svcName)
		}),
	}
	cmd.Flags().StringVarP(&push.DockerLogin.User, "docker-user", "u", "",
		"The docker login user, used before pushing image")
	cmd.Flags().StringVarP(&push.DockerLogin.Password, "docker-password", "p", "",
		"The docker login password, used before pushing image")
	return cmd
}

func imagesCmdFor(svcName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `images [env]`,
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
