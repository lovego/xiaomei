package elastic

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/lovego/xiaomei/utils/httputil"
)

func (es *ES) GetP(path string, bodyData map[string]interface{}, data interface{}) {
	es.Do(http.MethodGet, path, bodyData, data)
}

func (es *ES) PostP(path string, bodyData map[string]interface{}, data interface{}) {
	es.Do(http.MethodGet, path, bodyData, data)
}

func (es *ES) Do(method, path string, bodyData map[string]interface{}, data interface{}) {
	var body io.Reader
	if bodyData != nil {
		buf, err := json.Marshal(bodyData)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(buf)
	}
	resp := httputil.Do(http.MethodPut, es.Uri(path), nil, body)
	if resp != nil {
		resp.Body.Close()
	}
	resp.Ok().Json(data)
}
