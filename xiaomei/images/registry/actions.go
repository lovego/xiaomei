package registry

import (
	"fmt"
	"strings"

	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
)

func ListTimeTags(svcName, env string) {
	imgName := conf.GetService(svcName, env).ImageName()
	tags := []string{}
	for _, tag := range Tags(imgName) {
		if strings.HasPrefix(tag, env) && tag != env {
			tags = append(tags, tag[len(env):])
		}
	}
	fmt.Println(strings.Join(tags, "\n"))
}
