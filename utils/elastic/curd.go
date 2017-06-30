package elastic

import (
	"net/http"
	"net/url"

	"github.com/lovego/xiaomei/utils/httputil"
)

// 增
func (es *ES) Create(path string, bodyData, data interface{}) error {
	resp, err := httputil.Put(es.Uri(path), nil, bodyData)
	if err != nil {
		return err
	}
	if err := resp.Check(http.StatusOK, http.StatusCreated); err != nil {
		return err
	}
	return resp.Json(data)
}

// 删
func (es *ES) Delete(path string, data interface{}) error {
	return httputil.DeleteJson(es.Uri(path), nil, nil, data)
}

// 改
func (es *ES) Update(path string, bodyData, data interface{}) error {
	return httputil.PostJson(es.Uri(path+`/_update`), nil, bodyData, data)
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
	total int, data []map[string]interface{}, err error,
) {
	result := QueryResult{}
	uri, err := url.Parse(path + `/_search`)
	if err != nil {
		return
	}
	uri.Query().Set(`filter_path`, `hits.total,hits.hits._source`)

	if err = httputil.PostJson(es.Uri(uri.String()), nil, bodyData, &result); err != nil {
		return
	}
	total = result.Hits.Total
	for _, hit := range result.Hits.Hits {
		data = append(data, hit.Source)
	}
	return
}

func (es *ES) Ensure(path string, def interface{}) error {
	if ok, err := es.Exist(path); err != nil {
		return err
	} else if !ok {
		return es.Create(path, def, nil)
	}
	return nil
}

func (es *ES) Exist(path string) (bool, error) {
	resp, err := httputil.Head(es.Uri(path), nil, nil)
	if err != nil {
		return false, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, resp.CodeError()
	}
}
