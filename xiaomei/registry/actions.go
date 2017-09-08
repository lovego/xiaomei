package registry

import (
	"fmt"
	"strings"

	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
)

func ListTimeTags(svcName, env string) {
	imgName := conf.GetService(svcName, env).ImageName()
	tags := Tags(imgName)
	fmt.Println(strings.Join(tags, "\n"))
}
