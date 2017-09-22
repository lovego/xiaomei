package elastic

import (
	"encoding/json"
	"fmt"
	"log"
)

type BulkError interface {
	FailedItems() [][2]interface{}
	Error() string
}

type bulkError struct {
	typ         string
	inputs      [][2]interface{}
	results     []map[string]map[string]interface{}
	failedItems [][2]interface{}
}

func (b bulkError) FailedItems() [][2]interface{} {
	if b.failedItems == nil {
		failedItems := make([][2]interface{}, 0)
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

type BulkDeleteError interface {
	FailedItems() []interface{}
	Error() string
}

type bulkDeleteError struct {
	inputs      []interface{}
	results     []map[string]map[string]interface{}
	failedItems []interface{}
}

func (b bulkDeleteError) FailedItems() []interface{} {
	if b.failedItems == nil {
		failedItems := make([]interface{}, 0)
		for i, result := range b.results {
			res := result[`delete`]
			if res[`error`] != nil {
				failedItems = append(failedItems, b.inputs[i])
			}
		}
		b.failedItems = failedItems
	}
	return b.failedItems
}

func (b bulkDeleteError) Error() string {
	var errs []interface{}
	for _, result := range b.results {
		info := result[`delete`]
		if info[`error`] != nil {
			errs = append(errs, info)
		}
	}
	buf, err := json.MarshalIndent(errs, ``, `  `)
	if err != nil {
		log.Println(`marshal elastic bulk errors: `, err)
	}
	return fmt.Sprintf("bulk delete errors(%d of %d)\n%s\n", len(errs), len(b.inputs), buf)
}
