package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintJson(v interface{}) {
	if data, err := json.MarshalIndent(v, ``, `  `); err != nil {
		log.Panic(err)
	} else {
		fmt.Println(string(data))
	}
}
