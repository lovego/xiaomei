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
