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
	result, err := es.getQueryResult(path, `hits.total,hits.hits._source`, bodyData)
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
	result, err := es.getQueryResult(path, `hits.total,hits.hits._id,hits.hits._source`, bodyData)
	if err == nil {
		total = result.Hits.Total
		for _, hit := range result.Hits.Hits {
			hit.Source[`_id`] = hit.Id
			data = append(data, hit.Source)
		}
	}
	return
}

func (es *ES) QueryIds(path string, bodyData interface{}) (
	total int, data []string, err error,
) {
	result, err := es.getQueryResult(path, `hits.total,hits.hits._id`, bodyData)
	if err == nil {
		total = result.Hits.Total
		for _, hit := range result.Hits.Hits {
			data = append(data, hit.Id)
		}
	}
	return
}

func (es *ES) getQueryResult(path, filterPath string, bodyData interface{}) (*QueryResult, error) {
	uri, err := url.Parse(es.Uri(path))
	if err != nil {
		return nil, err
	}
	uri.Path += `/_search`
	q := uri.Query()
	q.Set(`filter_path`, filterPath)
	uri.RawQuery = q.Encode()

	result := &QueryResult{}
	if err = es.client.PostJson(uri.String(), nil, bodyData, result); err != nil {
		return nil, err
	}
	return result, nil
}
