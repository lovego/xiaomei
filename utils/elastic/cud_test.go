package elastic

import (
	"testing"
)

func TestCURD(t *testing.T) {
	createEmptyUsers()

	testES.Create(`/users/1`, map[string]interface{}{`name`: `lilei`, `age`: 21}, nil)
	testES.Create(`/users/2`, map[string]interface{}{`name`: `hanmeimei`, `age`: 19}, nil)
	testES.Create(`/users/3`, map[string]interface{}{`name`: `tom`, `age`: 22}, nil)

	testES.Delete(`/users/3`, nil)

	testES.Update(`/users/1`, map[string]map[string]int{`doc`: {`age`: 31}}, nil)
	testES.Update(`/users/2`, map[string]map[string]int{`doc`: {`age`: 29}}, nil)

	checkLiLeiAndHanMeiMei(t)
}
