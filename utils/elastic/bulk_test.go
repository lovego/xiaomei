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

	testES.BulkUpdate(`/users`, []map[string]interface{}{
		map[string]interface{}{`_id`: 1, `doc`: {`age`: 31}},
		map[string]interface{}{`_id`: 2, `doc`: {`age`: 29}},
	})

	checkLiLeiAndHanMeiMei(t)
}
