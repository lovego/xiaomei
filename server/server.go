package server

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/server/renderer"
	"github.com/bughou-go/xiaomei/server/renderer/funcs"
)

type Server struct {
	*Router
	Renderer       *renderer.Renderer
	FilterFunc     func(req *Request, res *Response) bool
	LayoutDataFunc func(
		layout string, data interface{}, req *Request, res *Response,
	) interface{}
}

func New(router *Router) *Server {
	return &Server{
		Router: router,
		Renderer: renderer.New(
			path.Join(config.Root(), `views`), `layout/default`,
			config.Data.Env != `dev`, funcs.Map(),
		),
	}
}

func (s *Server) ListenAndServe() {
	addr := config.CurrentAppServer().AppAddr + `:` + config.Data.AppPort

	fmt.Printf("%s listen at %s\n", time.Now().Format(`2006-01-02 15:04:05 -0700`), addr)

	if err := http.ListenAndServe(addr, http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			req := NewRequest(request)
			res := NewResponse(response, req, s.Renderer, s.LayoutDataFunc)

			defer handleError(time.Now(), req, res)

			// 如果返回true，继续交给路由处理
			if s.FilterFunc == nil || s.FilterFunc(req, res) {
				s.Router.Handle(req, res, notFound)
			}
		})); err != nil {
		panic(err)
	}
}
