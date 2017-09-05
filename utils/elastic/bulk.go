package elastic

import (
	"encoding/json"
	"log"
)

type BulkResult struct {
	Errors bool                                `json:"errors"`
	Items  []map[string]map[string]interface{} `json:"items"`
}

func (es *ES) BulkCreate(path string, data [][2]interface{}) error {
	if len(data) <= 0 {
		return nil
	}
	body, err := makeBulkCreate(data)
	if err != nil {
		return err
	}
	return es.BulkDo(path, body, `create`, data)
}

func (es *ES) BulkUpdate(path string, data [][2]interface{}) error {
	if len(data) <= 0 {
		return nil
	}
	return es.BulkDo(path, makeBulkUpdate(data), `update`, data)
}

func (es *ES) BulkDo(path string, body, typ string, data [][2]interface{}) error {
	result := BulkResult{}
	if err := es.client.PostJson(es.Uri(path+`/_bulk`), nil, body, &result); err != nil {
		return err
	}
	if !result.Errors {
		return nil
	}
	return bulkError{typ: typ, inputs: data, results: result.Items}
}

/*
	create es 格式：
	{ "create" : {"_id" : "2"} }
	{ "k": "v", ... }
	args:
	[ [ _id, data ], ...]
*/
func makeBulkCreate(rows [][2]interface{}) (result string, err error) {
	for _, row := range rows {
		id, data := row[0], row[1]
		if id == nil {
			if id, err = GenUUID(); err != nil {
				return ``, err
			}
		}

		meta, err := json.Marshal(map[string]map[string]interface{}{`create`: {`_id`: id}})
		if err != nil {
			log.Panic(err)
		}

		result += string(meta) + "\n" + makeBulkData(data) + "\n"
	}
	return
}

/*
	upsert es 格式：
	{ "update" : {"_id" : "2"} }
	{ "doc" : {"field" : "value"}, "upsert" : {"field" : "value"}, ... }
	args:
	[ [ _id, data ], ...]
*/
func makeBulkUpdate(rows [][2]interface{}) (result string) {
	for _, row := range rows {
		id, data := row[0], row[1]
		if id == nil {
			log.Panicf("must have _id: %v", row)
		}

		meta, err := json.Marshal(map[string]map[string]interface{}{`update`: {`_id`: id}})
		if err != nil {
			log.Panic(err)
		}

		result += string(meta) + "\n" + makeBulkData(data) + "\n"
	}
	return
}

func makeBulkData(data interface{}) string {
	switch data.(type) {
	case map[string]interface{}:
		content, err := json.Marshal(data.(map[string]interface{}))
		if err != nil {
			log.Panic(err)
		}
		return string(content)
	case []byte:
		return string(data.([]byte))
	default:
		log.Panic(`unexpected data type`)
	}
	return ``
}
