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
	result := BulkResult{}
	httputil.Post(es.Uri(path), nil, MakeBulkCreateBody(data)).Ok().Json(&result)
	if !result.Errors {
		return
	}
	var errs [][2]map[string]interface{}
	for i, item := range result.Items {
		info := item[`create`]
		if v, ok := info[`created`].(bool); !ok || !v {
			errs = append(errs, [2]map[string]interface{}{data[i], info})
		}
	}
	errsBuf, err := json.Marshal(errs)
	if err != nil {
		panic(err)
	}
	panic(fmt.Sprintf(`bulk create errors(%d of %d): %s`, len(errs), len(data), errsBuf))
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
