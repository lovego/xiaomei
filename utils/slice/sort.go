package slice

import (
	"reflect"
	"sort"
)

func Sort(data interface{}, key string) {
	sort.Sort(sortable{key, reflect.ValueOf(data)})
}

type sortable struct {
	sort  string
	value reflect.Value
}

func (d sortable) Len() int { return d.value.Len() }

func (d sortable) Less(i, j int) bool {
	iv := d.fieldValue(i, d.sort)
	jv := d.fieldValue(j, d.sort)
	if iv != nil && jv != nil {
		switch v := iv.(type) {
		case int:
			return v < jv.(int)
		default:
			return false
		}
	} else if iv == nil {
		return false
	} else if jv == nil {
		return true
	} else {
		return false
	}
}

func (d sortable) fieldValue(index int, name string) interface{} {
	if index >= d.value.Len() {
		return nil
	}
	v := d.value.Index(index)
	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(v.Interface())
	}
	switch v.Kind() {
	case reflect.Struct:
		return v.FieldByName(name).Interface()
	case reflect.Map:
		return v.MapIndex(reflect.ValueOf(name)).Interface()
	default:
		panic(`unsupported kind: ` + v.Kind().String())
	}
}

func (d sortable) Swap(i, j int) {
	iv := d.value.Index(i)
	jv := d.value.Index(j)
	tmp := iv.Interface()
	iv.Set(jv)
	jv.Set(reflect.ValueOf(tmp))
}
