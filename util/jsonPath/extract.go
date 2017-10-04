package jsonPath

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NodePrime/jsonpath"
)

// Extract extracts data from JSON bytes, matched by JSON path
func Extract(data []byte, path string) (string, error) {
	// Resolving paths
	paths, err := jsonpath.ParsePaths(path)
	if err != nil {
		return "", err
	}

	// Searching within provided data
	eval, err := jsonpath.EvalPathsInBytes(data, paths)
	if err != nil {
		return "", err
	}

	// Unmarshal JSON into interface
	var ii interface{}
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	if err = d.Decode(&ii); err != nil {
		return "", err
	}

	for {
		if result, ok := eval.Next(); ok {
			if value := seek(ii, result.Keys); len(value) > 0 {
				return value, nil
			}
		} else {
			break
		}
	}

	return "", nil
}

// seek searches for path value
func seek(data interface{}, keys []interface{}) string {
	if len(keys) == 0 {
		return fmt.Sprintf("%v", data) // Empty keys
	}
	k := keys[0]
	switch k.(type) {
	case int:
		// Slice
		index := k.(int)
		if slc, ok := data.([]interface{}); ok {
			if index >= len(slc) {
				return "" // Index out of range
			}

			return seek(slc[index], keys[1:])
		}
		return "" // Not a slice
	default:
		// Map
		mapKey := fmt.Sprintf("%q", k)
		mapKey = mapKey[1 : len(mapKey)-1]
		if mp, ok := data.(map[string]interface{}); ok {
			if mapValue, ok := mp[mapKey]; ok {
				return seek(mapValue, keys[1:])
			}
		}
		return ""
	}
}
