package elastic

import (
	"testing"
)

func TestBulk(t *testing.T) {
	createEmptyUsers()

	testES.BulkCreate(`/users`, []map[string]interface{}{
		{`_id`: 1, `name`: `lilei`, `age`: 21},
		{`_id`: 2, `name`: `hanmeimei`, `age`: 20},
		{`_id`: 3, `name`: `tom`, `age`: 22},
	})

	testES.Delete(`/users/3`, nil)

	testES.BulkUpdate(`/users`, [][2]interface{}{
		{1, map[string]map[string]int{`doc`: {`age`: 31}}},
		{2, map[string]map[string]int{`doc`: {`age`: 29}}},
	})

	checkLiLeiAndHanMeiMei(t)
}
