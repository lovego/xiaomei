package core

import (
	"github.com/lovego/goa"
	"{{ .ModulePath }}/core/goods"
	"{{ .ModulePath }}/core/orders"
	"{{ .ModulePath }}/core/stocks"
)

func Routes(router *goa.RouterGroup) {
	goods.Routes(router.Group("/goods", "商品"))
	orders.Routes(router.Group("/orders", "订单"))
	stocks.Routes(router.Group("/stocks", "库存"))
}
