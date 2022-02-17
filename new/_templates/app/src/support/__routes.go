package support

import (
	"github.com/lovego/goa"
	"{{ .ModulePath }}/support/accounts"
	"{{ .ModulePath }}/support/messages"
)

func Routes(router *goa.RouterGroup) {
	accounts.Routes(router.Group("/accounts", "账号"))
	messages.Routes(router.Group("/messages", "消息"))
}
