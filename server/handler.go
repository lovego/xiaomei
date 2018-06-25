package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	loggerPkg "github.com/lovego/logger"
	"github.com/lovego/tracer"
	"github.com/lovego/xiaomei"
)

func (s *Server) Handler() (handler http.Handler) {
	sysRoutes(s.Router)

	handler = s
	if s.HandleTimeout > 0 {
		handler = http.TimeoutHandler(handler, s.HandleTimeout,
			fmt.Sprintf(`ServeHTTP timeout after %s.`, s.HandleTimeout),
		)
	}
	return
}

func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	req := xiaomei.NewRequest(request, s.Session)
	res := xiaomei.NewResponse(response, req, s.Session, s.Renderer, s.LayoutDataFunc)

	debug := req.URL.Query()["_debug"] != nil
	logger.Record(debug, func(ctx context.Context) error {
		startTime := tracer.GetSpan(ctx).At
		psData.Add(request.Method, request.URL.Path, startTime)
		defer psData.Remove(request.Method, request.URL.Path, startTime)

		req.SetContext(ctx)
		if strings.HasPrefix(req.URL.Path, `/_`) || s.FilterFunc == nil || s.FilterFunc(req, res) {
			if !s.Router.Handle(req, res) {
				handleNotFound(res)
			}
		}
		return res.GetError()
	}, func() {
		handleServerError(res)
	}, func(fields *loggerPkg.Fields) {
		logFields(fields, req, res, debug)
	})
}

func handleNotFound(res *xiaomei.Response) {
	res.WriteHeader(404)
	if res.Size() <= 0 {
		res.Json(map[string]string{"code": "404", "message": "Not Found."})
	}
}

func handleServerError(res *xiaomei.Response) {
	res.WriteHeader(500)
	if res.Size() <= 0 {
		res.Json(map[string]string{"code": "server-err", "message": "Fatal Server Error."})
	}
}
