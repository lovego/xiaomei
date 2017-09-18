package elastic

import (
	"net/url"
)

// 查
type QueryResult struct {
	Hits QueryHits `json:"hits"`
}

type QueryHits struct {
	Total int        `json:"total"`
	Hits  []QueryHit `json:"hits"`
}

type QueryHit struct {
	Id     string                 `json:"_id"`
	Source map[string]interface{} `json:"_source"`
}

func (es *ES) Query(path string, bodyData interface{}) (
	total int, data []map[string]interface{}, err error,
) {
	result, err := es.getQueryResult(path, bodyData, false)
	if err == nil {
		total = result.Hits.Total
		for _, hit := range result.Hits.Hits {
			data = append(data, hit.Source)
		}
	}
	return
}

func (es *ES) QueryWithId(path string, bodyData interface{}) (
	total int, data []map[string]interface{}, err error,
) {
	result, err := es.getQueryResult(path, bodyData, true)
	total = result.Hits.Total
	for _, hit := range result.Hits.Hits {
		hit.Source[`_id`] = hit.Id
		data = append(data, hit.Source)
	}
	return
}

func (es *ES) getQueryResult(path string, bodyData interface{}, withId bool) (*QueryResult, error) {
	uri, err := url.Parse(path + `/_search`)
	if err != nil {
		return nil, err
	}
	if withId {
		uri.Query().Set(`filter_path`, `hits.total,hits.hits._id,hits.hits._source`)
	} else {
		uri.Query().Set(`filter_path`, `hits.total,hits.hits._source`)
	}

	result := &QueryResult{}
	if err = es.client.PostJson(es.Uri(uri.String()), nil, bodyData, result); err != nil {
		return nil, err
	}
	return result, nil
}
