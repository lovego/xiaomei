package elastic

import (
	"fmt"
	"testing"
)

var testES = New(`http://192.168.202.39:9200/bughou-test`)

func TestCURD(t *testing.T) {
	testCreateEmptyUsers()

	testES.Create(`/users/1`, map[string]interface{}{
		`name`: `lilei`, `age`: 21,
	}, nil)
	testES.Create(`/users/2`, map[string]interface{}{
		`name`: `hanmeimei`, `age`: 20,
	}, nil)
	testES.Create(`/users/3`, map[string]interface{}{
		`name`: `tom`, `age`: 17,
	}, nil)

	testES.Delete(`/users/3`, nil)

	testES.Update(`/users/1`, map[string]interface{}{`doc`: map[string]int{
		`age`: 31,
	}}, nil)
	testES.Update(`/users/2`, map[string]interface{}{`doc`: map[string]int{
		`age`: 30,
	}}, nil)

	var result map[string]interface{}
	testES.Get(`/_refresh`, nil, &result)
	// fmt.Println(result)

	fmt.Println(testES.Query(`/users`, nil))
}

func testCreateEmptyUsers() {
	testES.Delete(`/`, nil)

	testES.Ensure(`/`, nil)
	testES.Ensure(`/_mapping/users`, map[string]interface{}{
		"properties": map[string]interface{}{
			"name":       map[string]string{"type": "keyword"},
			"salt":       map[string]string{"type": "keyword"},
			"password":   map[string]string{"type": "keyword"},
			"created_at": map[string]string{"type": "date", "format": "yyyy-MM-dd'T'HH:mm:ssZ"},
		},
	})
}
