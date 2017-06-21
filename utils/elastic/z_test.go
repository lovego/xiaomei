package elastic

import (
	"testing"
)

var testES = New(`http://192.168.202.39:9200/bughou-test`)

func TestAll(t *testing.T) {
	testES.Ensure(`/`, nil)
	testES.Ensure(`/_mapping/users`, map[string]interface{}{
		"properties": map[string]interface{}{
			"name":       map[string]string{"type": "keyword"},
			"salt":       map[string]string{"type": "keyword"},
			"password":   map[string]string{"type": "keyword"},
			"created_at": map[string]string{"type": "date", "format": "yyyy-MM-dd'T'HH:mm:ssZ"},
		},
	})
	testES.Create(`/users/1`, map[string]interface{}{
		`name`: `lilei`, `age`: 18,
	}, nil)
	testES.Create(`/users/2`, map[string]interface{}{
		`name`: `hanmeimei`, `age`: 17,
	}, nil)
	testES.Create(`/users/3`, map[string]interface{}{
		`name`: `tom`, `age`: 17,
	}, nil)
	testES.Delete(`/users/3`, nil)
}
