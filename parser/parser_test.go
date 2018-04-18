package parser

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/walesey/dicelang/histogram"
	"github.com/walesey/dicelang/util"

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

func Test_Parser_Execute_EG1_Adds_To_1_0(t *testing.T) {
	parser := NewParser(bytes.NewReader([]byte("hist 16d6.4+.d6.4+.d6.4+.d6.4+")))
	result, err := parser.Execute()
	assert.Nil(t, err)
	var totalP float64
	for _, hc := range result.([]histogram.HistogramColumn) {
		totalP += hc.P
	}
	assert.EqualValues(t, 1.0, util.Round(totalP, .5, 5))
}

func Test_Parser_Execute_EG2_Adds_To_1(t *testing.T) {
	parser := NewParser(bytes.NewReader([]byte("hist 6d3.add.d6.3+.d6.3+.d6.not.4+.4")))
	result, err := parser.Execute()
	assert.Nil(t, err)
	var totalP float64
	for _, hc := range result.([]histogram.HistogramColumn) {
		totalP += hc.P
	}
	assert.EqualValues(t, 1.0, util.Round(totalP, .5, 5))
}
