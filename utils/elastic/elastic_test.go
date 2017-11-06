package elastic

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

var testES = New(`http://192.168.202.39:9200/bughou-test`)

func createEmptyUsers() {
	testES.Delete(`/`, nil)

	testES.Ensure(`/`, nil)
	testES.Ensure(`/_mapping/users`, map[string]interface{}{
		"properties": map[string]interface{}{
			"name": map[string]string{"type": "keyword"},
			"age":  map[string]string{"type": "integer"},
		},
	})
}

func checkLiLeiAndHanMeiMei(t *testing.T) {
	testES.Get(`/_refresh`, nil, nil)

	total, docs, err := testES.Query(`/users`, map[string]map[string]string{`sort`: {`age`: `desc`}})
	if err != nil {
		log.Panic(err)
	}
	expectTotal := 2
	expectDocs := []map[string]interface{}{
		{`name`: `lilei`, `age`: json.Number(`31`)},
		{`name`: `hanmeimei`, `age`: json.Number(`29`)},
	}
	if total != expectTotal {
		t.Errorf("expect total: %d, got: %d\n", expectTotal, total)
	}
	if !reflect.DeepEqual(docs, expectDocs) {
		t.Errorf(`
expect docs: %v
        got: %v
`, expectDocs, docs)
	}
}
