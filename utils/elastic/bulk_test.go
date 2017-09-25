package elastic

import (
	"testing"
)

func TestBulk(t *testing.T) {
	createEmptyUsers()

	testES.BulkCreate(`/users`, [][2]interface{}{
		{1, map[string]interface{}{`name`: `lilei`, `age`: 21}},
		{2, map[string]interface{}{`name`: `hanmeimei`, `age`: 20}},
		{3, map[string]interface{}{`name`: `tom`, `age`: 22}},
	})

	testES.Delete(`/users/3`, nil)

	testES.BulkUpdate(`/users`, [][2]interface{}{
		{1, map[string]map[string]interface{}{`doc`: {`age`: 31}}},
		{2, map[string]map[string]interface{}{`doc`: {`age`: 29}}},
	})

	checkLiLeiAndHanMeiMei(t)
}
