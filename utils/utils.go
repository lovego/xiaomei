package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintJson(v interface{}) {
	data, err := json.MarshalIndent(v, ``, `  `)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(data))
}

func SprintJson(v interface{}) string {
	data, err := json.MarshalIndent(v, ``, `  `)
	if err != nil {
		log.Panic(err)
	}
	return string(data)
}
