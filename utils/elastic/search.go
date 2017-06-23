package elastic

import (
	"github.com/lovego/xiaomei/utils/httputil"
)

type SearchResult struct {
	Took    int          `json:"took"`
	Timeout bool         `json:"time_out"`
	Shards  SearchShards `json:"_shards"`
	Hits    SearchHits   `json:"hits"`
}

type SearchShards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

type SearchHits struct {
	Total int         `json:"total"`
	Hits  []SearchHit `json:"hits"`
}

type SearchHit struct {
	Source map[string]interface{} `json:"_source"`
}

func (es *ES) Search(path string, bodyData map[string]interface{}) SearchResult {
	result := SearchResult{}
	httputil.Post(es.Uri(path+`/_search`), nil, bodyData).Ok().Json(&result)
	return result
}
