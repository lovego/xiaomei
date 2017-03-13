package xm

import (
	// "github.com/googollee/go-socket.io"
	"regexp"
)

// 字符串路由处理函数
type StrRouteHandler func(*Request, *Response)

// 正则路由处理函数，第三个参数是正则匹配的结果
type RegRouteHandler func(*Request, *Response, []string)

type RegRoute struct {
	reg     *regexp.Regexp
	handler RegRouteHandler
}

type Router struct {
	// 基础路径
	basePath string
	// 字符串路由 method     path
	strRoutes map[string]map[string]StrRouteHandler
	// 正则路由   method     base_path
	regRoutes map[string]map[string][]RegRoute
	// SocketIO  *socketio.Server
}

func NewRouter() *Router {
	return &Router{
		strRoutes: make(map[string]map[string]StrRouteHandler),
		regRoutes: make(map[string]map[string][]RegRoute),
	}
}

/*
func (r *Router) On(e string, f interface{}) {
	if err := r.SocketIO.Of(r.basePath).On(e, f); err != nil {
		panic(err)
	}
}

func (r *Router) Of(p string) socketio.Namespace {
	return r.SocketIO.Of(p)
}
*/
