package release

import (
	"errors"
	"fmt"

	"github.com/lovego/xiaomei/utils/slice"
	"github.com/spf13/cobra"
)

type cmdFunc func(c *cobra.Command, args []string) error

func NoArgCall(work func() error) cmdFunc {
	return func(c *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.New(`redundant args.`)
		}
		return work()
	}
}

func Arg1Call(arg1 string, work func(string) error) cmdFunc {
	return func(c *cobra.Command, args []string) error {
		switch len(args) {
		case 0:
			if arg1 == `` {
				return errors.New(`one argument is required.`)
			}
		case 1:
			arg1 = args[0]
		default:
			return errors.New(`more than one arguments given.`)
		}
		return work(arg1)
	}
}

func EnvCall(work func(string) error) cmdFunc {
	return func(c *cobra.Command, args []string) error {
		var env string
		switch len(args) {
		case 0:
			env = `dev`
		case 1:
			env = args[0]
			if !slice.ContainsString(getEnvs(), env) {
				return fmt.Errorf("env %s not defined in cluster.yml", env)
			}
		default:
			return errors.New(`more than one arguments given.`)
		}
		return work(env)
	}
}

func Env1Call(work func(string, string) error) cmdFunc {
	return func(c *cobra.Command, args []string) error {
		var env, arg1 string
		switch len(args) {
		case 0:
			env = `dev`
		case 1:
			env = args[0]
		case 2:
			env = args[0]
			arg1 = args[1]
		default:
			return errors.New(`more than two arguments given.`)
		}
		if !slice.ContainsString(getEnvs(), env) {
			return fmt.Errorf("env %s not defined in cluster.yml", env)
		}
		return work(env, arg1)
	}
}
