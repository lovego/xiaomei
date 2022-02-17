package generic

import (
	"github.com/lovego/goa"
	"{{ .ModulePath }}/generic/im"
	"{{ .ModulePath }}/generic/files"
)

func Routes(router *goa.RouterGroup) {
	im.Routes(router.Group("/goods", "实时通知"))
	files.Routes(router.Group("/orders", "文件上传、下载"))
}
