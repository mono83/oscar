package jsonpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	assert := assert.New(t)
	source := `{"id": 10, "data": [1, "bar", {"id": 8}]}`

	if result, err := Extract([]byte(source), "$.id"); assert.NoError(err) {
		assert.Equal("10", result)
	}

	if result, err := Extract([]byte(source), `$[id]`); assert.NoError(err) {
		assert.Equal("10", result)
	}

	if result, err := Extract([]byte(source), `$["xxx"]`); assert.NoError(err) {
		assert.Equal("", result)
	}

	if result, err := Extract([]byte(source), `$.xxx`); assert.NoError(err) {
		assert.Equal("", result)
	}

	if result, err := Extract([]byte(source), "$.data[0]"); assert.NoError(err) {
		assert.Equal("1", result)
	}

	if result, err := Extract([]byte(source), "$.data[1]"); assert.NoError(err) {
		assert.Equal("bar", result)
	}

	if result, err := Extract([]byte(source), "$.data[2].id"); assert.NoError(err) {
		assert.Equal("8", result)
	}

	if result, err := Extract([]byte(source), "$.data[3]"); assert.NoError(err) {
		assert.Equal("", result)
	}
}

func TestExtractLongInt(t *testing.T) {
	assert := assert.New(t)
	source := `{"id": 1060000000001000003}`

	if result, err := Extract([]byte(source), "$.id"); assert.NoError(err) {
		assert.Equal("1060000000001000003", result)
	}
}

func TestExtractLongFloat(t *testing.T) {
	assert := assert.New(t)
	source := `{"id": 1060000000001000003.00001}`

	if result, err := Extract([]byte(source), "$.id"); assert.NoError(err) {
		assert.Equal("1060000000001000003.00001", result)
	}
}
