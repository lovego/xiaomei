package strmap

import (
	"reflect"
	"testing"
)

var sm = StrMap(map[string]interface{}{
	`key1`: `value1`,
	`key2`: map[interface{}]interface{}{`key`: `value`},
	`key3`: []interface{}{`value1`, `value2`, `value3`},
})

func TestGet(t *testing.T) {
	if got := sm.Get(`nonexist`); got != nil {
		t.Errorf("expect: %v got: %v\n", nil, got)
	}

	if got := sm.Get(`key1`); got != nil {
		t.Errorf("expect: %v got: %v\n", nil, got)
	}

	if got := sm.Get(`key1`).Get(`nonexist1`).Get(`nonexist2`); got != nil {
		t.Errorf("expect: %v got: %v\n", nil, got)
	}

	if got := sm.Get(`key3`); got != nil {
		t.Errorf("expect: %v got: %v\n", nil, got)
	}

	expect := StrMap(map[string]interface{}{`key`: `value`})
	if got := sm.Get(`key2`); !reflect.DeepEqual(got, expect) {
		t.Errorf("expect: %v got: %v\n", expect, got)
	}
}

func TestGetString(t *testing.T) {
	if expect, got := ``, sm.GetString(`nonexist`); got != expect {
		t.Errorf("expect: %v got: %v\n", expect, got)
	}

	if expect, got := `value1`, sm.GetString(`key1`); got != expect {
		t.Errorf("expect: %v got: %v\n", expect, got)
	}

	if expect, got := ``, sm.GetString(`key2`); got != expect {
		t.Errorf("expect: %v got: %v\n", expect, got)
	}

	if expect, got := ``, sm.GetString(`key3`); got != expect {
		t.Errorf("expect: %v got: %v\n", expect, got)
	}
}

func TestGetStringSlice(t *testing.T) {
	if got := sm.GetStringSlice(`nonexist`); got != nil {
		t.Errorf("expect: %v got: %v\n", nil, got)
	}

	if expect, got := []string{`value1`}, sm.GetStringSlice(`key1`); !reflect.DeepEqual(got, expect) {
		t.Errorf("expect: %v got: %v\n", nil, got)
	}

	if got := sm.GetStringSlice(`key2`); got != nil {
		t.Errorf("expect: %v got: %v\n", nil, got)
	}

	expect, got := []string{`value1`, `value2`, `value3`}, sm.GetStringSlice(`key3`)
	if !reflect.DeepEqual(got, expect) {
		t.Errorf("expect: %v got: %v\n", expect, got)
	}
}
