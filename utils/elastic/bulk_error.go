package elastic

import (
	"encoding/json"
	"fmt"
	"log"
)

type BulkError interface {
	Items() [][2]interface{}
	Error() string
}

type bulkError struct {
	typ         string
	inputs      []map[string]interface{}
	results     []map[string]map[string]interface{}
	failedItems []map[string]interface{}
}

func (b bulkError) FailedItems() []map[string]interface{} {
	if b.failedItems == nil {
		failedItems := make([]map[string]interface{}, 0)
		for i, result := range b.results {
			res := result[b.typ]
			if res[`error`] != nil {
				failedItems = append(failedItems, b.inputs[i])
			}
		}
		b.failedItems = failedItems
	}
	return b.failedItems
}

func (b bulkError) Error() string {
	var errs []interface{}
	for _, result := range b.results {
		info := result[b.typ]
		if info[`error`] != nil {
			errs = append(errs, info)
		}
	}
	buf, err := json.MarshalIndent(errs, ``, `  `)
	if err != nil {
		log.Println(`marshal elastic bulk errors: `, err)
	}
	return fmt.Sprintf("bulk %s errors(%d of %d)\n%s\n",
		b.typ, len(errs), len(b.inputs), buf,
	)
}
