package elastic

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/lovego/xiaomei/utils/httputil"
)

type SearchResult struct {
	Hits struct {
		Total int `json:"total"`
		Hits  []struct {
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (es *ES) Get(path string, bodyData map[string]interface{}) (
	total int, data []map[string]interface{},
) {
	return es.Search(http.MethodGet, path, bodyData)
}

func (es *ES) Post(path string, bodyData map[string]interface{}) (
	total int, data []map[string]interface{},
) {
	return es.Search(http.MethodPost, path, bodyData)
}

func (es *ES) Search(method, path string, bodyData map[string]interface{}) (
	total int, data []map[string]interface{},
) {
	result := SearchResult{}
	es.Do(method, path, bodyData, &result)
	total = result.Hits.Total
	for _, hit := range result.Hits.Hits {
		data = append(data, hit.Source)
	}
	return
}

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
