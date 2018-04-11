package dice

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/walesey/dicelang/histogram"
)

func Test_MultiDice_Hist(t *testing.T) {
	d := MultiDice{Dice{Size: 6}, 2} // 2d6.add
	h := d.Hist()
	expected := map[int]float64{
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
	}
	histogram.RoundHistogram(expected)
	histogram.RoundHistogram(h)
	assert.EqualValues(t, expected, h)
}
