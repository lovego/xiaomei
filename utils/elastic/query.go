package elastic

import (
	"net/url"
)

func (es *ES) Query(path string, bodyData interface{}) (
	total int, data []map[string]interface{}, err error,
) {
	uri, err := url.Parse(path)
	if err != nil {
		return 0, nil, err
	}
	q := uri.Query()
	q.Set(`filter_path`, `hits.total,hits.hits._source`)
	uri.RawQuery = q.Encode()

	result, err := es.Search(uri.String(), bodyData)
	if err == nil {
		total = result.Hits.Total
		err = result.Hits.Sources(&data)
	}
	return
}

func (es *ES) QueryWithId(path string, bodyData interface{}, key string) (
	total int, data []map[string]interface{}, err error,
) {
	uri, err := url.Parse(path)
	if err != nil {
		return 0, nil, err
	}
	q := uri.Query()
	q.Set(`filter_path`, `hits.total,hits.hits._source,hits.hits._id`)
	uri.RawQuery = q.Encode()

	result, err := es.Search(uri.String(), bodyData)
	if err == nil {
		total = result.Hits.Total
		err = result.Hits.SourcesWithId(&data, key)
	}
	return
}
