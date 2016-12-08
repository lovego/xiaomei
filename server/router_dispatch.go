package server

import "strings"

// 处理请求
func (r *Router) Handle(req *Request, res *Response) bool {
	method := strings.ToUpper(req.Method)
	path := cleanPath(req.URL.Path)
	if r.strRoutesMatch(method, path, req, res) || r.regRoutesMatch(method, path, req, res) {
		return true
	}
	return false
}

func (r *Router) strRoutesMatch(method string, path string, req *Request, res *Response) bool {
	routes := r.strRoutes[method]
	if routes == nil {
		return false
	}
	handler := routes[path]
	if handler == nil {
		return false
	}
	handler(req, res)
	return true
}

// 按斜线作为分隔符从深到浅依次匹配
func (r *Router) regRoutesMatch(method string, path string, req *Request, res *Response) bool {
	routes := r.regRoutes[method]
	if routes == nil {
		return false
	}
	for p := path; p != ``; {
		// 上一层路径
		if i := strings.LastIndexByte(p, '/'); i > 0 {
			p = p[:i]
		} else {
			p = ``
		}
		if slice := routes[p]; slice != nil {
			mp := path[len(p):]
			for _, route := range slice {
				if m := route.reg.FindStringSubmatch(mp); m != nil {
					route.handler(req, res, m)
					return true
				}
			}
		}
	}
	return false
}
