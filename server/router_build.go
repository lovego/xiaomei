package xiaomei

import (
	"path"
	"regexp"
)

// Group 提供带basePath的路由，代码更简洁，正则匹配更高效。
// p只能是字符串路径，不能是正则表达式。
func (r *Router) Group(p string) *Router {
	p = cleanPath(p)
	if r.basePath != `` {
		p = path.Join(r.basePath, p)
	}
	return &Router{
		basePath:  p,
		strRoutes: r.strRoutes,
		regRoutes: r.regRoutes,
		// SocketIO:  r.SocketIO,
	}
}

func (r *Router) Get(p string, handler StrRouteHandler) *Router {
	return r.addStrRoute(`GET`, p, handler)
}

func (r *Router) Post(p string, handler StrRouteHandler) *Router {
	return r.addStrRoute(`POST`, p, handler)
}

func (r *Router) GetPost(p string, handler StrRouteHandler) *Router {
	return r.addStrRoute(`GET`, p, handler).addStrRoute(`POST`, p, handler)
}

func (r *Router) Put(p string, handler StrRouteHandler) *Router {
	return r.addStrRoute(`PUT`, p, handler)
}

func (r *Router) Delete(p string, handler StrRouteHandler) *Router {
	return r.addStrRoute(`DELETE`, p, handler)
}

func (r *Router) GetX(reg string, handler RegRouteHandler) *Router {
	return r.addRegRoute(`GET`, reg, handler)
}

func (r *Router) PostX(reg string, handler RegRouteHandler) *Router {
	return r.addRegRoute(`POST`, reg, handler)
}

func (r *Router) GetPostX(reg string, handler RegRouteHandler) *Router {
	return r.addRegRoute(`GET`, reg, handler).addRegRoute(`POST`, reg, handler)
}

func (r *Router) PutX(reg string, handler RegRouteHandler) *Router {
	return r.addRegRoute(`PUT`, reg, handler)
}

func (r *Router) DeleteX(reg string, handler RegRouteHandler) *Router {
	return r.addRegRoute(`DELETE`, reg, handler)
}

// 增加字符串路由
func (r *Router) addStrRoute(method string, p string, handler StrRouteHandler) *Router {
	p = cleanPath(p)
	if r.basePath != `` {
		p = path.Join(r.basePath, p)
	}
	if r.strRoutes[method] == nil {
		r.strRoutes[method] = make(map[string]StrRouteHandler)
	}
	r.strRoutes[method][p] = handler
	return r
}

// 增加正则路由
func (r *Router) addRegRoute(method string, reg string, handler RegRouteHandler) *Router {
	if r.regRoutes[method] == nil {
		r.regRoutes[method] = make(map[string][]RegRoute)
	}
	r.regRoutes[method][r.basePath] = append(r.regRoutes[method][r.basePath], RegRoute{
		regexp.MustCompile(`^` + reg + `$`), handler,
	})
	return r
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	return path.Clean(p)
}
