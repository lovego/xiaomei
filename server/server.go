package server

import (
	"net"
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/server/funcs"
	"github.com/bughou-go/xiaomei/server/xm"
	"github.com/bughou-go/xiaomei/server/xm/renderer"
	"github.com/bughou-go/xiaomei/server/xm/session"
	"github.com/fatih/color"
)

func init() {
	if n := runtime.NumCPU() - 1; n >= 1 {
		runtime.GOMAXPROCS(n)
	}
}

type Server struct {
	FilterFunc     func(req *xm.Request, res *xm.Response) bool
	Router         *xm.Router
	Session        session.Session
	Renderer       *renderer.Renderer
	LayoutDataFunc func(layout string, data interface{}, req *xm.Request, res *xm.Response) interface{}
}

func NewSession() session.Session {
	return session.NewCookieSession(http.Cookie{
		Name: config.Name(),
		Path: `/`,
	}, config.Secret())
}

func NewRenderer() *renderer.Renderer {
	return renderer.New(
		path.Join(config.Root(), `views`), `layout/default`, config.Env() != `dev`, funcs.Map(),
	)
}

func (s *Server) ListenAndServe() {
	const addr = `:3000`
	ln := listen(addr)
	config.Log(color.GreenString(`started. (` + addr + `)`))

	svr := http.Server{Handler: s.Handler()}
	if err := svr.Serve(ln); err != nil {
		panic(err)
	}
}

func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			req := xm.NewRequest(request, s.Session)
			res := xm.NewResponse(response, req, s.Session, s.Renderer, s.LayoutDataFunc)

			var notFound bool
			defer handleError(time.Now(), req, res, &notFound)

			// 如果返回true，继续交给路由处理
			if s.FilterFunc == nil || s.FilterFunc(req, res) {
				notFound = !s.Router.Handle(req, res)
			}
		})
}

func listen(addr string) net.Listener {
	if ln, err := net.Listen(`tcp`, addr); err != nil {
		panic(err)
	} else {
		return tcpKeepAliveListener{ln.(*net.TCPListener)}
	}
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
