package develop

import (
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Build() error {
	if err := build(); err != nil {
		return err
	}
	if err := Spec(``); err != nil {
		return err
	}
	Assets(nil)
	cmd.Run(cmd.O{}, `git`, `status`)
	return nil
}
