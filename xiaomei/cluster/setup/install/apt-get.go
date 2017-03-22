package install

import (
	"os"
	"time"

	"github.com/lovego/xiaomei/utils/cmd"
)

var aptGetNeedUpdate = true

func AptGet(pkg string) error {
	if aptGetNeedUpdate {
		// 超过3天没更新源，则更新
		if fi, err := os.Stat(`/var/lib/apt/periodic/update-success-stamp`); err != nil {
			return err
		} else if time.Since(fi.ModTime()) > 72*time.Hour {
			cmd.Run(cmd.O{}, `sudo`, `apt-get`, `update`, `--fix-missing`)
		}
		aptGetNeedUpdate = false
	}
	_, err := cmd.Run(cmd.O{}, `sudo`, `apt-get`, `install`, `-y`, pkg)
	return err
}
