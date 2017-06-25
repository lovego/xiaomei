package elastic

import (
	"net/http"
	"net/url"

	"github.com/lovego/xiaomei/utils/httputil"
)

// 增
func (es *ES) Create(path string, bodyData, data interface{}) {
	resp := httputil.Put(es.Uri(path), nil, bodyData)

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
	default:
		panic(`unexpected response: ` + resp.Status + "\n" + string(resp.GetBody()))
	}
	resp.Json(data)
}

// 删
func (es *ES) Delete(path string, data interface{}) {
	httputil.Delete(es.Uri(path), nil, nil).Ok().Json(data)
}

// 改
func (es *ES) Update(path string, bodyData, data interface{}) {
	httputil.Post(es.Uri(path+`/_update`), nil, bodyData).Ok().Json(data)
}

// 查
type QueryResult struct {
	Hits QueryHits `json:"hits"`
}

type QueryHits struct {
	Total int        `json:"total"`
	Hits  []QueryHit `json:"hits"`
}

type QueryHit struct {
	Source map[string]interface{} `json:"_source"`
}

func (es *ES) Query(path string, bodyData interface{}) (
	total int, data []map[string]interface{},
) {
	result := QueryResult{}
	uri, err := url.Parse(path + `/_search`)
	uri.Query().Set(`filter_path`, `hits.total,hits.hits._source`)
	if err != nil {
		panic(err)
	}
	httputil.Post(es.Uri(uri.String()), nil, bodyData).Ok().JsonUseNumber(&result)
	total = result.Hits.Total
	for _, hit := range result.Hits.Hits {
		data = append(data, hit.Source)
	}
	return
}

func (es *ES) Ensure(path string, def interface{}) {
	if !es.Exist(path) {
		es.Create(path, def, nil)
	}
}

func (es *ES) Exist(path string) bool {
	resp := httputil.Head(es.Uri(path), nil, nil)
	if resp != nil {
		defer resp.Body.Close()
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
