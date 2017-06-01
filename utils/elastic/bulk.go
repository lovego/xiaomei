package elastic

import (
	"encoding/json"
	"strings"

	"github.com/lovego/xiaomei/utils/httputil"
	"github.com/nu7hatch/gouuid"
)

type BulkResult struct {
	Errors bool                                `json:"errors"`
	Items  []map[string]map[string]interface{} `json:"items"`
}

func (es *ES) BulkCreate(path string, data []map[string]interface{}) {
	result := BulkResult{}
	httputil.Post(es.Uri(path), nil, MakeBulkCreateBody(data)).Ok().Json(&result)
	if !result.Errors {
		return
	}
	for i, item := range result.Items {
	}
}

func MakeBulkCreateBody(rows []map[string]interface{}) (result string) {
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
