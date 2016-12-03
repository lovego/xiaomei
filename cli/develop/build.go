package develop

import (
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Build() {
	if !build() {
		return
	}

	if !Spec(``) {
		return
	}

	Assets(nil)

	cmd.Run(cmd.O{}, `git`, `status`)
}
