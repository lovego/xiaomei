package elastic

import (
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
	httputil.Do(method, es.Uri(path), nil, bodyData).Ok().Json(data)
}
