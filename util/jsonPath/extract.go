package jsonPath

import (
	"encoding/json"
	"fmt"
	"github.com/yalp/jsonpath"
	"bytes"
)

// Extract extracts data from JSON bytes, matched by JSON path
func Extract(data []byte, path string) (string, error) {
	var iface interface{}

	d := json.NewDecoder(bytes.NewBuffer(data))
	d.UseNumber()

	if err := d.Decode(&iface); err != nil {
		return "", err
	}

	response, err := jsonpath.Read(iface, path)
	if err != nil {
		return "", filterError(err)
	}

	return fmt.Sprint(response), nil
}

// filterError filters errors, emitted by yalp/jsonpath library
func filterError(err error) error {
	msg := err.Error()
	var str string
	var num int
	if n, _ := fmt.Sscanf(msg, "out of bound array access at %d", &num); n == 1 {
		return nil
	} else if n, _ := fmt.Sscanf(msg, "no key '%s' for object at %d", &str, &num); n == 1 {
		return nil
	} else if n, _ := fmt.Sscanf(msg, "child '%s' not found in JSON object at %d", &str, &num); n == 1 {
		return nil
	}

	return err
}
