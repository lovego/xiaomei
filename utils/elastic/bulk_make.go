package elastic

import (
	"encoding/json"
	"log"
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

/*
	delete es 格式：
	{ "delete" : {  "_id" : "2" } }
	args:
	[ [ _id, data ], ...]
*/
func makeBulkDelete(rows []interface{}) (result string) {
	for _, id := range rows {
		if id == nil {
			log.Panicf("must have _id: %v", id)
		}
		meta, err := json.Marshal(map[string]map[string]interface{}{`update`: {`_id`: id}})
		if err != nil {
			log.Panic(err)
		}
		result += string(meta) + "\n"
	}
	return
}
