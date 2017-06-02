package elastic

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lovego/xiaomei/utils/httputil"
	"github.com/nu7hatch/gouuid"
)

type BulkResult struct {
	Errors bool                                `json:"errors"`
	Items  []map[string]map[string]interface{} `json:"items"`
}

func (es *ES) BulkCreate(path string, data []map[string]interface{}) {
	es.BulkDo(path, MakeBulkCreate(data), `create`, `created`, data)
}

func (es *ES) BulkUpdate(path string, data []map[string]interface{}) {
	es.BulkDo(path, MakeBulkUpdate(data), `update`, `updated`, data)
}

func (es *ES) BulkDo(path string, body, typ, expect string, data []map[string]interface{}) {
	result := BulkResult{}
	httputil.Post(es.Uri(path), nil, body).Ok().Json(&result)
	if !result.Errors {
		return
	}
	var errs [][2]map[string]interface{}
	for i, item := range result.Items {
		info := item[typ]
		if v, ok := info[`result`].(string); !ok || v != expect {
			errs = append(errs, [2]map[string]interface{}{data[i], info})
		}
	}
	errsBuf, err := json.Marshal(errs)
	if err != nil {
		panic(err)
	}
	panic(fmt.Sprintf(`bulk %s errors(%d of %d): %s`, typ, len(errs), len(data), errsBuf))
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

func GenUUID() string {
	if uid, err := uuid.NewV4(); err != nil {
		panic(err)
	} else {
		return strings.Replace(uid.String(), `-`, ``, -1)
	}
}

/*
	upsert es 格式：
	{ "update" : {"_id" : "2"} }
	{ "doc" : {"field" : "value"}, "doc_as_upsert" : true }
*/
func MakeBulkUpdate(rows []map[string]interface{}) (result string) {
	for _, row := range rows {
		id := row[`_id`]
		if id == nil {
			rowBuf, err := json.Marshal(row)
			if err != nil {
				panic(err)
			}
			panic(`bulk update no _id for: ` + string(rowBuf))
		}
		meta, err := json.Marshal(map[string]interface{}{`update`: map[string]interface{}{`_id`: id}})
		if err != nil {
			panic(err)
		}

		delete(row, `_id`)
		content, err := json.Marshal(map[string]interface{}{"doc": row, "doc_as_upsert": true})
		if err != nil {
			panic(err)
		}
		result += string(meta) + "\n" + string(content) + "\n"
	}
	return
}
