package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parser_SimpleExec(t *testing.T) {
	parser := NewParser(bytes.NewReader([]byte("hist 4d6.4+")))
	expected := `{"0": 0.0625, "1": 0.25, "2": 0.375, "3": 0.25, "4": 0.0625}`
	result, err := parser.Execute()
	assert.Nil(t, err)
	assert.JSONEq(t, expected, result)
}
