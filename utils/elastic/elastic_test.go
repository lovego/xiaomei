package elastic

import (
	"encoding/json"
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

	result, err := testES.Search(`/users`, map[string]map[string]string{`sort`: {`age`: `desc`}})
	if err != nil {
		t.Fatal(err)
	}

	expectTotal := 2
	if result.Hits.Total != expectTotal {
		t.Errorf("expect total: %d, got: %d\n", expectTotal, result.Hits.Total)
	}

	expectDocs := []map[string]interface{}{
		{`name`: `lilei`, `age`: json.Number(`31`)},
		{`name`: `hanmeimei`, `age`: json.Number(`29`)},
	}
	var docs []map[string]interface{}
	if err := result.Hits.Sources(&docs); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(docs, expectDocs) {
		t.Errorf(`
expect docs: %v
        got: %v
`, expectDocs, docs)
	}
}
