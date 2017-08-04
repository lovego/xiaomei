package elastic

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type BulkError interface {
	Items() [][2]interface{}
	Error() string
}

type bulkError struct {
	typ     string
	inputs  reflect.Value
	results []map[string]map[string]interface{}
	items   [][2]interface{}
}

func (b bulkError) Items() [][2]interface{} {
	if b.items == nil {
		items := make([][2]interface{}, 0)
		for i, result := range b.results {
			res := result[b.typ]
			if res[`error`] != nil {
				items = append(items, [2]interface{}{b.inputs.Index(i).Interface(), res})
			}
		}
		b.items = items
	}
	return b.items
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
		b.typ, len(errs), b.inputs.Len(), buf,
	)
}
