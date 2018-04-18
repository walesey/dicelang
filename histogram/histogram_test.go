package histogram

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/walesey/dicelang/util"
)

type D2 struct{}

func (d D2) Resolve() int { return rand.Intn(2) + 1 }

func (d D2) Hist() map[int]float64 { return map[int]float64{1: 0.5, 2: 0.5} }

type D4 struct{}

func (d D4) Resolve() int { return rand.Intn(4) + 1 }

func (d D4) Hist() map[int]float64 { return map[int]float64{1: 0.25, 2: 0.25, 3: 0.25, 4: 0.25} }

func Test_Aggregate(t *testing.T) {
	h := Aggregate(D4{}, D4{}).Hist()
	expected := map[int]float64{
		2: 1.0 / 16.0,
		3: 2.0 / 16.0,
		4: 3.0 / 16.0,
		5: 4.0 / 16.0,
		6: 3.0 / 16.0,
		7: 2.0 / 16.0,
		8: 1.0 / 16.0,
	}
	assert.EqualValues(t, expected, h)
}

func Test_Multiply(t *testing.T) {
	h := Multiply(D2{}, D4{}).Hist()
	expected := map[int]float64{
		1: 4.0 / 32.0,
		2: 5.0 / 32.0,
		3: 6.0 / 32.0,
		4: 7.0 / 32.0,
		5: 4.0 / 32.0,
		6: 3.0 / 32.0,
		7: 2.0 / 32.0,
		8: 1.0 / 32.0,
	}
	assert.EqualValues(t, expected, h)
}

func Test_Multiply_Adds_To_One(t *testing.T) {
	h := Multiply(D2{}, D2{}, D2{}, D2{}, D2{}, D2{}).Hist()
	var totalP float64
	for _, p := range h {
		totalP += p
	}
	assert.EqualValues(t, 1.0, util.Round(totalP, .5, 5))
}
