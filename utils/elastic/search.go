package elastic

import (
	"encoding/json"
	"net/url"
)

type SearchResult struct {
	Took         int             `json:"took"`
	Timeout      bool            `json:"time_out"`
	Shards       Shards          `json:"_shards"`
	Hits         SearchHits      `json:"hits"`
	Aggregations json.RawMessage `json:"aggregations"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

func (es *ES) Search(path string, bodyData interface{}) (*SearchResult, error) {
	uri, err := url.Parse(es.Uri(path))
	if err != nil {
		return nil, err
	}
	uri.Path += `/_search`
	result := &SearchResult{}
	if err := es.client.PostJson(uri.String(), nil, bodyData, result); err != nil {
		return nil, err
	}
	return result, nil
}
