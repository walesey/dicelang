package dice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Dice_Hist(t *testing.T) {
	d := Dice{Size: 6} // d6
	h := d.Hist()
	oneSixth := 1.0 / 6.0
	expected := map[int]float64{1: oneSixth, 2: oneSixth, 3: oneSixth, 4: oneSixth, 5: oneSixth, 6: oneSixth}
	assert.EqualValues(t, expected, h)
}

func Test_Dice_Hist_Weighted(t *testing.T) {
	d := Dice{Size: 3, Weights: []float64{0.2, 0.3, 0.5}} // d6
	h := d.Hist()
	expected := map[int]float64{1: 0.2, 2: 0.3, 3: 0.5}
	assert.EqualValues(t, expected, h)
}

func Test_Dice_Values(t *testing.T) {
	d := Dice{Size: 6} // d6
	v := d.Values()
	expected := []int{1, 2, 3, 4, 5, 6}
	assert.EqualValues(t, expected, v)
}

func Test_Dice_Values_D4(t *testing.T) {
	d := Dice{Size: 4} // d4
	v := d.Values()
	expected := []int{1, 2, 3, 4}
	assert.EqualValues(t, expected, v)
}
