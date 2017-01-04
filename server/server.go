package server

import (
	"net/http"
	"path"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/server/funcs"
	"github.com/bughou-go/xiaomei/server/xm"
	"github.com/bughou-go/xiaomei/server/xm/renderer"
	"github.com/bughou-go/xiaomei/server/xm/session"
)

type Server struct {
	FilterFunc     func(req *xm.Request, res *xm.Response) bool
	Router         *xm.Router
	Session        session.Session
	Renderer       *renderer.Renderer
	LayoutDataFunc func(layout string, data interface{}, req *xm.Request, res *xm.Response) interface{}
}

func NewSession() session.Session {
	return session.NewCookieSession(http.Cookie{
		Name: config.App.Name(),
		Path: `/`,
	}, config.App.Secret())
}

func NewRenderer() *renderer.Renderer {
	return renderer.New(
		path.Join(config.App.Root(), `views`), `layout/default`, config.App.Env() != `dev`, funcs.Map(),
	)
}

func (s *Server) ListenAndServe() {
	addr := config.Servers.CurrentAppServer().AppAddr()

	if err := http.ListenAndServe(addr, http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			req := xm.NewRequest(request, s.Session)
			res := xm.NewResponse(response, req, s.Session, s.Renderer, s.LayoutDataFunc)

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
