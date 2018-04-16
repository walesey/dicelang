package parser

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parser_Execute(t *testing.T) {
	parser := NewParser(bytes.NewReader([]byte("hist 4d6.4+")))
	expected := `[{"V": 0, "P": 0.0625}, {"V": 1, "P": 0.25}, {"V": 2, "P": 0.375}, {"V": 3, "P": 0.25}, {"V": 4, "P": 0.0625}]`
	result, err := parser.Execute()
	assert.Nil(t, err)
	jsonResult, err := json.Marshal(result)
	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(jsonResult))
}

func Test_Parser_Execute_Aggregate(t *testing.T) {
	parser := NewParser(bytes.NewReader([]byte("hist [2d6.4+, d6.add]")))
	expected := `[{"V":1,"P":0.0417},{"V":2,"P":0.125},{"V":3,"P":0.1667},{"V":4,"P":0.1667},{"V":5,"P":0.1667},{"V":6,"P":0.1667},{"V":7,"P":0.125},{"V":8,"P":0.0417}]`
	result, err := parser.Execute()
	assert.Nil(t, err)
	jsonResult, err := json.Marshal(result)
	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(jsonResult))
}
