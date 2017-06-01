package elastic

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/lovego/xiaomei/utils/httputil"
)

type ES struct {
	BaseAddr string
}

func New(addr string) *ES {
	return &ES{BaseAddr: strings.TrimSuffix(addr, `/`)}
}

func (es *ES) Ensure(path string, def map[string]interface{}) {
	if !es.Exist(path) {
		es.Create(path, def)
	}
}

func (es *ES) Exist(path string) bool {
	resp := httputil.Do(http.MethodHead, es.Uri(path), nil, nil)
	if resp != nil {
		resp.Body.Close()
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return true
	case http.StatusNotFound:
		return false
	default:
		panic(`unexpected response: ` + resp.Status + "\n" + string(resp.GetBody()))
	}
}

func (es *ES) Create(path string, def map[string]interface{}) {
	var body io.Reader
	if def != nil {
		buf, err := json.Marshal(def)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(buf)
	}
	resp := httputil.Do(http.MethodPut, es.Uri(path), nil, body)
	if resp != nil {
		resp.Body.Close()
	}
	resp.Ok()
}

func (es *ES) Uri(path string) string {
	if strings.HasPrefix(path, `/`) {
		return es.BaseAddr + path
	} else {
		return es.BaseAddr + `/` + path
	}
}
