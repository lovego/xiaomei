package elastic

import (
	"encoding/json"
)

/*
	upsert es 格式：
	{ "update" : {"_id" : "2"} }
	{ "doc" : {"field" : "value"}, "upsert" : {"field" : "value"} }
*/
func MakeBulkUpdate(rows [][2]interface{}) (result string) {
	for _, row := range rows {
		id, updateDef := row[0], row[1]
		meta, err := json.Marshal(map[string]map[string]interface{}{`update`: {`_id`: id}})
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
