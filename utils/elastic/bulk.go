package elastic

import (
	"reflect"

	"github.com/lovego/xiaomei/utils/httputil"
)

type BulkResult struct {
	Errors bool                                `json:"errors"`
	Items  []map[string]map[string]interface{} `json:"items"`
}

func (es *ES) BulkCreate(path string, data []map[string]interface{}) error {
	return es.BulkDo(path, MakeBulkCreate(data), `create`, data)
}

func (es *ES) BulkUpdate(path string, data [][2]interface{}) error {
	return es.BulkDo(path, MakeBulkUpdate(data), `update`, data)
}

func (es *ES) BulkDo(path string, body, typ string, data interface{}) error {
	result := BulkResult{}
	if err := httputil.PostJson(es.Uri(path+`/_bulk`), nil, body, &result); err != nil {
		return err
	}
	if !result.Errors {
		return nil
	}
	return bulkError{typ: typ, inputs: reflect.ValueOf(data), results: result.Items}
}
