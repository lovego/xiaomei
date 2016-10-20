package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJson(v interface{}) {
	data, err := json.MarshalIndent(v, ``, `  `)
	fmt.Println(string(data), err)
}
