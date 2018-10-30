package release

import (
	"errors"

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
		var arg0 string
		switch len(args) {
		case 0:
		case 1:
			arg0 = args[0]
		default:
			return errors.New(`more than one arguments given.`)
		}
		if env, err := CheckEnv(arg0); err == nil {
			return work(env)
		} else {
			return err
		}
	}
}

func Env1Call(work func(string, string) error) cmdFunc {
	return func(c *cobra.Command, args []string) error {
		var arg0, arg1 string
		switch len(args) {
		case 0:
		case 1:
			arg0 = args[0]
		case 2:
			arg0, arg1 = args[0], args[1]
		default:
			return errors.New(`more than two arguments given.`)
		}
		if env, err := CheckEnv(arg0); err == nil {
			return work(env, arg1)
		} else {
			return err
		}
	}
}

func EnvSliceCall(work func(string, []string) error) cmdFunc {
	return func(c *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New(`at leaset two arguments required.`)
		}
		if env, err := CheckEnv(args[0]); err == nil {
			return work(env, args[1:])
		} else {
			return err
		}
	}
}
