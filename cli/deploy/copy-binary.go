package deploy

import (
	"path"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func copyBinary(sshAddr, tag string) {
	binPath := path.Join(config.Deploy.Path(), `release/bins`, tag)
	if config.IsLocalEnv() {
		if cmd.Fail(cmd.O{}, `test`, `-x`, binPath) {
			cmd.Run(cmd.O{Panic: true}, `cp`, config.App.Bin(), binPath)
		}
	} else {
		if cmd.Fail(cmd.O{}, `ssh`, sshAddr, `test`, `-x`, binPath) {
			cmd.Run(cmd.O{Panic: true}, `scp`, config.App.Bin(), sshAddr+`:`+binPath)
		}
	}
}
