package server

import (
	"net"
	"net/http"
	"os"
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

func (s *Server) ListenAndServe() {
	if s.Server == nil {
		s.Server = &http.Server{}
	}
	s.Server.Handler = s.Handler()

	port := os.Getenv(`GOPORT`)
	if port == `` {
		port = `3000`
	}
	addr := `:` + port
	listener := listen(addr)
	utils.Log(color.GreenString(`started.(` + addr + `)`))

	if err := s.Server.Serve(listener); err != nil {
		panic(err)
	}
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
