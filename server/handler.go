package server

import (
	"time"

	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/server/xm"
)

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

func handleError(t time.Time, req *xm.Request, res *xm.Response, notFound *bool) {
	if *notFound {
		handleNotFound(req, res)
	}

	err := recover()
	if err != nil {
		handleServerError(req, res)
	}
	if err == nil && req.URL.Path == alivePath {
		return
	}
	log := writeLog(req, res, t, err)
	if err != nil {
		go config.Alarm(`500错误`, string(log))
	}
}

func handleNotFound(req *xm.Request, res *xm.Response) {
	res.WriteHeader(404)
	if res.Size() <= 0 {
		res.Json(map[string]string{"code": "404", "message": "Not Found."})
	}
}

func handleServerError(req *xm.Request, res *xm.Response) {
	res.WriteHeader(500)
	if res.Size() <= 0 {
		res.Json(map[string]string{"code": "500", "message": "Application Server Error."})
	}
}
