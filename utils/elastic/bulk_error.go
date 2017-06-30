package elastic

import (
	"bytes"
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
	errsCount := 0
	errsMap := make(map[string]int)
	for _, result := range b.results {
		if err := result[b.typ][`error`]; err != nil {
			if errInfo, ok := err.(map[string]interface{}); ok {
				if typ, ok := errInfo[`type`].(string); ok {
					errsCount++
					errsMap[typ]++
				}
			}
		}
	}
	buf := bytes.NewBufferString(fmt.Sprintf(
		"bulk %s errors(%d of %d)\n", b.typ, errsCount, b.inputs.Len(),
	))
	for typ, count := range errsMap {
		fmt.Fprintf(buf, "%s: %d\n", typ, count)
	}
	return buf.String()
}
