package dice

import (
	"testing"

	"github.com/walesey/dicelang/histogram"
	"github.com/walesey/dicelang/util"

	"github.com/stretchr/testify/assert"
)

func Test_DicePool_Hist(t *testing.T) {
	d6 := Dice{Size: 6}
	d := DicePool{Count: 4, Dice: d6, GTE: 4} // 4d6.4+
	h := d.Hist()
	expected := map[int]float64{0: 0.0625, 1: 0.25, 2: 0.375, 3: 0.25, 4: 0.0625}
	assert.EqualValues(t, expected, h)
}

func Test_DicePool_Multiply_Adds_To_One(t *testing.T) {
	d6 := Dice{Size: 6}
	dp := DicePool{Count: 16, Dice: d6, GTE: 4}
	dp2 := DicePool{Count: 1, Dice: d6, GTE: 4}
	h := histogram.Multiply(dp, dp2, dp2, dp2).Hist() // 16d6.4+.d6.4+.d6.4+.d6.4+

	var totalP float64
	for _, p := range h {
		totalP += p
	}
	assert.EqualValues(t, 1.0, util.Round(totalP, .5, 5))
}
