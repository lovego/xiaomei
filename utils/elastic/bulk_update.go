package elastic

import (
	"encoding/json"
	"log"
)

/*
	upsert es 格式：
	{ "update" : {"_id" : "2"} }
	{ "_id": id, "doc" : {"field" : "value"}, "upsert" : {"field" : "value"}, ... }
*/
func MakeBulkUpdate(rows []map[string]interface{}) (result string) {
	for _, row := range rows {
		id := row[`_id`]
		if id == nil {
			log.Panicf("must have _id: %v", row)
		}
		row = copyWithoutId(row)

		meta, err := json.Marshal(map[string]map[string]interface{}{`update`: {`_id`: id}})
		if err != nil {
			log.Panic(err)
		}

		content, err := json.Marshal(row)
		if err != nil {
			log.Panic(err)
		}
		result += string(meta) + "\n" + string(content) + "\n"
	}
	return
}
