package dice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DicePool_Hist(t *testing.T) {
	d6 := Dice{Size: 6}
	d := DicePool{Count: 4, Dice: d6, GTE: 4} // 4d6.4+
	h := d.Hist()
	expected := map[int]float64{0: 0.0625, 1: 0.25, 2: 0.375, 3: 0.25, 4: 0.0625}
	assert.EqualValues(t, expected, h)
}
