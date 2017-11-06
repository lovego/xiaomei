package elastic

import (
	"bytes"
	"encoding/json"
	"reflect"
)

type SearchHits struct {
	Total int         `json:"total"`
	Hits  []SearchHit `json:"hits"`
}

type SearchHit struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type"`
	Id     string          `json:"_id"`
	Score  string          `json:"_score"`
	Source json.RawMessage `json:"_source"`
}

func (h SearchHits) Sources(data interface{}) error {
	buf := bytes.NewBufferString(`[`)
	last := len(h.Hits) - 1
	for i, hit := range h.Hits {
		buf.Write(hit.Source)
		if i < last {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')

	decoder := json.NewDecoder(buf)
	decoder.UseNumber()
	return decoder.Decode(&data)
}

func (h SearchHits) Ids(path string, bodyData interface{}) (data []string) {
	for _, hit := range h.Hits {
		data = append(data, hit.Id)
	}
	return
}

func (h SearchHits) SourcesWithId(data interface{}, key string) error {
	err := h.Sources(data)
	if err != nil {
		h.SetIds(data, key)
	}
	return err
}

func (h SearchHits) SetIds(data interface{}, key string) {
	value := reflect.ValueOf(data)
	if kind := value.Kind(); kind == reflect.Ptr || kind == reflect.Interface {
		value = value.Elem()
	}
	length := value.Len() // should be an array or slice
	for i := 0; i < length; i++ {
		v := value.Index(i)
		if kind := v.Kind(); kind == reflect.Ptr || kind == reflect.Interface {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Struct:
			v.FieldByName(key).SetString(h.Hits[i].Id)
		case reflect.Map:
			v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(h.Hits[i].Id))
		}
	}
}
