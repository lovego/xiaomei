package elastic

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/lovego/xiaomei/utils/httputil"
)

type BulkResult struct {
	Errors bool                                `json:"errors"`
	Items  []map[string]map[string]interface{} `json:"items"`
}

func (es *ES) BulkCreate(path string, data []map[string]interface{}) {
	es.BulkDo(path, MakeBulkCreate(data), `create`, data)
}

func (es *ES) BulkUpdate(path string, data [][2]interface{}) {
	es.BulkDo(path, MakeBulkUpdate(data), `update`, data)
}

func (es *ES) BulkDo(path string, body, typ string, data interface{}) {
	result := BulkResult{}
	httputil.Post(es.Uri(path+`/_bulk`), nil, body).Ok().Json(&result)
	if !result.Errors {
		return
	}
	dataV := reflect.ValueOf(data)
	var errs [][2]interface{}
	for i, item := range result.Items {
		resp := item[typ]
		if resp[`error`] != nil {
			errs = append(errs, [2]interface{}{dataV.Index(i).Interface(), resp})
		}
	}
	errsBuf, err := json.Marshal(errs)
	if err != nil {
		panic(err)
	}
	panic(fmt.Sprintf(`bulk %s errors(%d of %d): %s`, typ, len(errs), dataV.Len(), errsBuf))
}

/*
	create es 格式：
	{ "create" : {"_id" : "2"} }
	{ "k": "v", ... }
*/
func MakeBulkCreate(rows []map[string]interface{}) (result string) {
	for _, row := range rows {
		meta, err := json.Marshal(map[string]interface{}{`create`: map[string]string{`_id`: GenUUID()}})
		if err != nil {
			panic(err)
		}

		content, err := json.Marshal(row)
		if err != nil {
			panic(err)
		}
		result += string(meta) + "\n" + string(content) + "\n"
	}
	return
}

/*
	upsert es 格式：
	{ "update" : {"_id" : "2"} }
	{ "doc" : {"field" : "value"}, "upsert" : {"field" : "value"} }
*/
func MakeBulkUpdate(rows [][2]interface{}) (result string) {
	for _, row := range rows {
		id, updateDef := row[0], row[1]
		meta, err := json.Marshal(map[string]interface{}{`update`: map[string]interface{}{`_id`: id}})
		if err != nil {
			panic(err)
		}

		content, err := json.Marshal(updateDef)
		if err != nil {
			panic(err)
		}
		result += string(meta) + "\n" + string(content) + "\n"
	}
	return
}
