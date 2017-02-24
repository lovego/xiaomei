package app

import (
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Build() error {
	if err := buildBinary(); err != nil {
		return err
	}
	if err := Spec(``); err != nil {
		return err
	}
	// Assets(nil)
	// build image
	cmd.Run(cmd.O{}, `git`, `status`)
	return nil
}
