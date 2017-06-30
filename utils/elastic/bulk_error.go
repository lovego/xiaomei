package elastic

import (
	"encoding/json"
	"fmt"
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
	buf, err := json.Marshal(b.Items())
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`bulk %s errors(%d of %d): %s`, b.typ, len(b.Items()), b.inputs.Len(), buf)
}
