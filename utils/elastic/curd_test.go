package elastic

import (
	"encoding/json"
	"reflect"
	"testing"
)

var testES = New(`http://192.168.202.39:9200/bughou-test`)

func TestCURD(t *testing.T) {
	testCreateEmptyUsers()

	testES.Create(`/users/1`, map[string]interface{}{`name`: `lilei`, `age`: 21}, nil)
	testES.Create(`/users/2`, map[string]interface{}{`name`: `hanmeimei`, `age`: 20}, nil)
	testES.Create(`/users/3`, map[string]interface{}{`name`: `tom`, `age`: 17}, nil)

	testES.Delete(`/users/3`, nil)

	testES.Update(`/users/1`, map[string]map[string]int{`doc`: {`age`: 31}}, nil)
	testES.Update(`/users/2`, map[string]map[string]int{`doc`: {`age`: 30}}, nil)

	testES.Get(`/_refresh`, nil, nil)

	total, docs := testES.Query(`/users`, map[string]map[string]string{`sort`: {`age`: `desc`}})
	expectTotal := 2
	expectDocs := []map[string]interface{}{
		{`name`: `lilei`, `age`: json.Number(`31`)},
		{`name`: `hanmeimei`, `age`: json.Number(`30`)},
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

func testCreateEmptyUsers() {
	testES.Delete(`/`, nil)

	testES.Ensure(`/`, nil)
	testES.Ensure(`/_mapping/users`, map[string]interface{}{
		"properties": map[string]interface{}{
			"name":       map[string]string{"type": "keyword"},
			"age":        map[string]string{"type": "integer"},
			"salt":       map[string]string{"type": "keyword"},
			"password":   map[string]string{"type": "keyword"},
			"created_at": map[string]string{"type": "date", "format": "yyyy-MM-dd'T'HH:mm:ssZ"},
		},
	})
}
