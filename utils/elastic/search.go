package elastic

import (
	"net/http"
	"net/url"
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
	uri, err := url.Parse(path)
	uri.Query().Set(`filter_path`, `hits.total,hits.hits._source`)
	if err != nil {
		panic(err)
	}
	es.Do(method, uri.String(), bodyData, &result)
	total = result.Hits.Total
	for _, hit := range result.Hits.Hits {
		data = append(data, hit.Source)
	}
	return
}
