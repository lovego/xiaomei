package server

import (
	"github.com/lovego/xiaomei/server/xm"
)

const alivePath = `/_alive`

func sysRoutes(router *xm.Router) {
	// 存活检测
	router.Root(`GET`, alivePath, func(req *xm.Request, res *xm.Response) {
		res.Write([]byte(`ok`))
	})
}
