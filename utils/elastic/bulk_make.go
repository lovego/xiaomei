package elastic

import (
	"encoding/json"
	"fmt"
)

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
			return ``, err
		}
		dataStr, err := makeBulkData(data)
		if err != nil {
			return ``, err
		}

		result += string(meta) + "\n" + dataStr + "\n"
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
func makeBulkUpdate(rows [][2]interface{}) (result string, err error) {
	for _, row := range rows {
		id, data := row[0], row[1]
		if id == nil {
			return ``, fmt.Errorf("empty _id: %v", row)
		}

		meta, err := json.Marshal(map[string]map[string]interface{}{`update`: {`_id`: id}})
		if err != nil {
			return ``, err
		}
		dataStr, err := makeBulkData(data)
		if err != nil {
			return ``, err
		}

		result += string(meta) + "\n" + dataStr + "\n"
	}
	return
}

func makeBulkData(data interface{}) (string, error) {
	switch value := data.(type) {
	case string:
		return value, nil
	case []byte:
		return string(value), nil
	default:
		content, err := json.Marshal(value)
		if err != nil {
			return ``, err
		}
		return string(content), nil
	}
}

/*
	delete es 格式：
	{ "delete" : {  "_id" : "2" } }
	args:
	[ [ _id, data ], ...]
*/
func makeBulkDelete(ids []string) (result string, err error) {
	for _, id := range ids {
		if id == `` {
			return ``, fmt.Errorf("has empty _id: %v", ids)
		}
		meta, err := json.Marshal(map[string]map[string]string{`delete`: {`_id`: id}})
		if err != nil {
			return ``, err
		}
		result += string(meta) + "\n"
	}
	return
}
