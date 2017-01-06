package deploy

import (
	"path"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func copyBinaries(sshAddr, tag string) {
	copyFmwkBin(sshAddr, tag)
	copyProjectBin(sshAddr, tag)
}

func copyFmwkBin(sshAddr, tag string) {
	fmwkBin, _ := cmd.Run(cmd.O{Output: true, Panic: true}, `which`, `xiaomei`)
	if config.IsLocalEnv() {
		return
	}
	if cmd.Fail(cmd.O{}, `ssh`, sshAddr, `which`, `xiaomei`) {
		cmd.Run(cmd.O{Panic: true}, `scp`, fmwkBin, sshAddr+`:/tmp/xiaomei`)
		cmd.Run(cmd.O{Panic: true}, `ssh`, sshAddr, `sudo mv /tmp/xiaomei /usr/local/bin`)
	}
}

func copyProjectBin(sshAddr, tag string) {
	tagBin := path.Join(config.Deploy.Path(), `release/bins`, tag)
	if config.IsLocalEnv() {
		if cmd.Fail(cmd.O{}, `test`, `-x`, tagBin) {
			cmd.Run(cmd.O{Panic: true}, `cp`, config.App.Bin(), tagBin)
		}
	} else {
		if cmd.Fail(cmd.O{}, `ssh`, sshAddr, `test`, `-x`, tagBin) {
			cmd.Run(cmd.O{Panic: true}, `scp`, config.App.Bin(), sshAddr+`:`+tagBin)
		}
	}
}
