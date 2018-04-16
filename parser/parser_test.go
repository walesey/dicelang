package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parser_SimpleExec(t *testing.T) {
	parser := NewParser(bytes.NewReader([]byte("hist 4d6.4+")))
	expected := `[{"V": 0, "P": 0.0625}, {"V": 1, "P": 0.25}, {"V": 2, "P": 0.375}, {"V": 3, "P": 0.25}, {"V": 4, "P": 0.0625}]`
	result, err := parser.Execute()
	assert.Nil(t, err)
	assert.JSONEq(t, expected, result)
}
