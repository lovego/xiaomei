package server

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/server/renderer"
	"github.com/bughou-go/xiaomei/server/renderer/funcs"
	"github.com/bughou-go/xiaomei/server/session"
)

type Server struct {
	FilterFunc     func(req *Request, res *Response) bool
	Router         *Router
	Session        session.Session
	Renderer       *renderer.Renderer
	LayoutDataFunc func(layout string, data interface{}, req *Request, res *Response) interface{}
}

func NewSession() session.Session {
	return session.NewCookieSession(http.Cookie{
		Name: config.App.Name(),
	}, config.App.Secret())
}

func NewRenderer() *renderer.Renderer {
	return renderer.New(
		path.Join(config.App.Root(), `views`), `layout/default`, config.App.Env() != `dev`, funcs.Map(),
	)
}

func (s *Server) ListenAndServe() {
	addr := config.Servers.Current().AppAddr + `:` + config.App.Port()

	fmt.Printf("%s listen at %s\n", time.Now().Format(`2006-01-02 15:04:05 -0700`), addr)

	if err := http.ListenAndServe(addr, http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			req := NewRequest(request, s.Session)
			res := NewResponse(response, req, s.Session, s.Renderer, s.LayoutDataFunc)

			var notFound bool
			defer handleError(time.Now(), req, res, &notFound)

			// 如果返回true，继续交给路由处理
			if s.FilterFunc == nil || s.FilterFunc(req, res) {
				notFound = !s.Router.Handle(req, res)
			}
		})); err != nil {
		panic(err)
	}
}
