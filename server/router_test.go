package xiaomei

// go test -c -gcflags "-N -l"

import (
	// `fmt`
	"net/http"
	"strings"
	"testing"
)

func notFound(req *Request, res *Response) {}

type testRouteData struct {
	method string
	path   string
	name   string
}

var testRoutes = []testRouteData{
	testRouteData{`get`, `/`, `index`},
	testRouteData{`get`, `/new`, `new`},
	testRouteData{`post`, `/`, `create`},
	testRouteData{`getx`, `/(\d+)`, `show`},
	testRouteData{`getx`, `/(\d+)/edit`, `edit`},
	testRouteData{`postx`, `/(\d+)/update`, `update`},
	testRouteData{`postx`, `/(\d+)/destroy`, `destroy`},
}

func TestNewRouter(t *testing.T) {
	var r *Router = NewRouter()
	if r == nil || r.strRoutes == nil || r.regRoutes == nil {
		t.Error()
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	matched := make(map[string]bool)
	for _, route := range testRoutes {
		testAddRoute(r, route, matched, t)
		testHandleReq(r, route, matched, ``, t)
	}
}

func TestGroupRouter(t *testing.T) {
	var prefix = `/admin`

	r := NewRouter()
	g := r.Group(prefix)
	matched := make(map[string]bool)
	for _, route := range testRoutes {
		testAddRoute(g, route, matched, t)
		testHandleReq(r, route, matched, prefix, t)
	}
}

func testAddRoute(r *Router, route testRouteData, matched map[string]bool, testing *testing.T) {
	switch route.method {
	case `get`:
		r.Get(route.path, func(req *Request, res *Response) {
			matched[route.name] = true
		})
	case `post`:
		r.Post(route.path, func(req *Request, res *Response) {
			matched[route.name] = true
		})
	case `getx`:
		r.GetX(route.path, func(req *Request, res *Response, params []string) {
			matched[route.name] = true
		})
	case `postx`:
		r.PostX(route.path, func(req *Request, res *Response, params []string) {
			matched[route.name] = true
		})
	}
}

func testHandleReq(
	r *Router, route testRouteData, matched map[string]bool, prefix string, t *testing.T,
) {
	method := strings.Replace(route.method, `x`, ``, 1)
	path := strings.Replace(route.path, `(\d+)`, `123`, 1)
	req, _ := http.NewRequest(method, prefix+path, nil)
	r.Handle(&Request{Request: req}, nil, notFound)
	if !matched[route.name] {
		t.Errorf(`%s not matched`, prefix+path)
	}
}
