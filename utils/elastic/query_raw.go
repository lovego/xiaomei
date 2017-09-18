package elastic

import (
	"bytes"
	"encoding/json"
	"net/url"
	"reflect"
)

// 查
type RawQueryResult struct {
	Hits RawQueryHits `json:"hits"`
}

type RawQueryHits struct {
	Total int           `json:"total"`
	Hits  []RawQueryHit `json:"hits"`
}

type RawQueryHit struct {
	Id     string          `json:"_id"`
	Source json.RawMessage `json:"_source"`
}

func (es *ES) QueryRaw(path string, bodyData interface{}) (
	total int, data []RawQueryHit, err error,
) {
	result, err := es.getRawQueryResult(path, bodyData, true)
	if err == nil {
		total = result.Hits.Total
		data = result.Hits.Hits
	}
	return
}

func (es *ES) QueryData(path string, bodyData interface{}, data interface{}) (
	total int, err error,
) {
	result, err := es.getRawQueryResult(path, bodyData, true)
	if err == nil {
		total = result.Hits.Total
		err = parseSource(result.Hits.Hits, data)
	}

	return
}

func (es *ES) QueryDataWithId(path string, bodyData interface{}, data interface{}) (
	total int, err error,
) {
	result, err := es.getRawQueryResult(path, bodyData, true)
	if err == nil {
		total = result.Hits.Total
		err = parseSource(result.Hits.Hits, data)
		if err == nil {
			setupDataIds(data, result.Hits.Hits)
		}
	}

	return
}

func (es *ES) getRawQueryResult(path string, bodyData interface{}, withId bool) (*RawQueryResult, error) {
	uri, err := url.Parse(path + `/_search`)
	if err != nil {
		return nil, err
	}
	if withId {
		uri.Query().Set(`filter_path`, `hits.total,hits.hits._id,hits.hits._source`)
	} else {
		uri.Query().Set(`filter_path`, `hits.total,hits.hits._source`)
	}

	result := &RawQueryResult{}
	if err = es.client.PostJson(es.Uri(uri.String()), nil, bodyData, result); err != nil {
		return nil, err
	}
	return result, nil
}

func parseSource(hits []RawQueryHit, data interface{}) error {
	source := bytes.NewBufferString(`[`)
	last := len(hits) - 1
	for i, hit := range hits {
		source.Write(hit.Source)
		if i < last {
			source.WriteByte(',')
		}
	}
	source.WriteByte(']')

	return json.Unmarshal(source.Bytes(), data)
}

func setupDataIds(data interface{}, hits []RawQueryHit) {
	value := reflect.ValueOf(data)
	if kind := value.Kind(); kind == reflect.Interface || kind == reflect.Ptr {
		value = value.Elem()
	}
	length := value.Len()
	for i := 0; i < length; i++ {
		v := value.Index(i)
		if kind := v.Kind(); kind == reflect.Interface || kind == reflect.Ptr {
			v = v.Elem()
		}
		v.FieldByName(`Id`).SetString(hits[i].Id)
	}
}
