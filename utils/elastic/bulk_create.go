package elastic

import (
	"encoding/json"
	"log"
)

/*
	create es 格式：
	{ "create" : {"_id" : "2"} }
	{ "_id": id, "k": "v", ... }
*/
func MakeBulkCreate(rows []map[string]interface{}) (result string, err error) {
	for _, row := range rows {
		id := row[`_id`]
		if id == nil {
			if id, err = GenUUID(); err != nil {
				return ``, err
			}
		} else {
			row = copyWithoutId(row)
		}

		meta, err := json.Marshal(map[string]map[string]interface{}{`create`: {`_id`: id}})
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

func copyWithoutId(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		if k != `_id` {
			result[k] = v
		}
	}
	return result
}
