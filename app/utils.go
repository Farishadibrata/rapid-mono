package app

import "encoding/json"

func Stringify(data map[string]interface{}) string {
	toBeMarshalled := make(map[string]interface{}, len(data))
	for k, v := range data {
		toBeMarshalled[k] = v
	}
	bytes, _ := json.Marshal(toBeMarshalled)

	return string(bytes)
}
