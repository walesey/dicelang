package dice

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/walesey/dicelang/histogram"
	"github.com/walesey/dicelang/util"
)

func Test_MultiDice_Hist(t *testing.T) {
	d := MultiDice{Dice{Size: 6}, 2} // 2d6.add
	expectedHist := histogram.Fixed(map[int]float64{
		2:  1.0 / 36.0,
		3:  2.0 / 36.0,
		4:  3.0 / 36.0,
		5:  4.0 / 36.0,
		6:  5.0 / 36.0,
		7:  6.0 / 36.0,
		8:  5.0 / 36.0,
		9:  4.0 / 36.0,
		10: 3.0 / 36.0,
		11: 2.0 / 36.0,
		12: 1.0 / 36.0,
	})
	expected := histogram.RoundHistogram(expectedHist, 5)
	h := histogram.RoundHistogram(d, 5)
	assert.EqualValues(t, expected, h)
}

func Test_MultiDice_Multiply_Adds_To_One(t *testing.T) {
	d6 := Dice{Size: 6}
	md := MultiDice{Count: 2, Dice: d6}
	h := histogram.Multiply(md, md, md).Hist() // 2d6.add.2d6.add.2d6.add

	var totalP float64
	for _, p := range h {
		totalP += p
	}

	assert.EqualValues(t, 1.0, util.Round(totalP, .5, 5))
}
