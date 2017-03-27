package server

import (
	"fmt"
	"net"
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/server/funcs"
	"github.com/lovego/xiaomei/server/xm"
	"github.com/lovego/xiaomei/server/xm/renderer"
	"github.com/lovego/xiaomei/server/xm/session"
	"github.com/lovego/xiaomei/utils"
)

func init() {
	if n := runtime.NumCPU() - 1; n >= 1 {
		runtime.GOMAXPROCS(n)
	}
}

type Server struct {
	*http.Server
	HandleTimeout  time.Duration
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
		path.Join(config.Root(), `views`), `layout/default`, !config.DevMode(), funcs.Map(),
	)
}

const alivePath = `/_alive`

func (s *Server) ListenAndServe() {
	s.Router.Root(`GET`, alivePath, func(req *xm.Request, res *xm.Response) {
		res.Write([]byte(`ok`))
	})

	if s.Server == nil {
		s.Server = &http.Server{}
	}
	s.Server.Handler = s.Handler()

	const addr = `:3000`
	listener := listen(addr)
	utils.Log(color.GreenString(`started. (` + addr + `)`))

	if err := s.Server.Serve(listener); err != nil {
		panic(err)
	}
}

func (s *Server) Handler() (handler http.Handler) {
	handler = http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			req := xm.NewRequest(request, s.Session)
			res := xm.NewResponse(response, req, s.Session, s.Renderer, s.LayoutDataFunc)

			var notFound bool
			defer handleError(time.Now(), req, res, &notFound)

			// 如果返回true，继续交给路由处理
			if req.Request.URL.Path == alivePath || s.FilterFunc == nil || s.FilterFunc(req, res) {
				notFound = !s.Router.Handle(req, res)
			}
		})
	if s.HandleTimeout > 0 {
		handler = http.TimeoutHandler(handler, s.HandleTimeout,
			fmt.Sprintf(`ServeHTTP timeout after %s.`, s.HandleTimeout),
		)
	}
	return
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
