package strnum

import (
	"encoding/json"
)

func JsonNumber2float64(data []map[string]interface{}) error {
	for _, row := range data {
		for k, v := range row {
			if num, ok := v.(json.Number); ok {
				if f64, err := num.Float64(); err == nil {
					row[k] = f64
				} else {
					return err
				}
			}
		}
	}
	return nil
}
