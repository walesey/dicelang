package dice

import "github.com/walesey/dicelang/histogram"

// DicePool - a pool of dice rolls with a GTE value
type DicePool struct {
	Dice  Dice
	Count int
	GTE   int
}

func (dp DicePool) Resolve() int {
	var total int
	for i := 0; i < dp.Count; i++ {
		roll := dp.Dice.Resolve()
		if roll >= dp.GTE {
			total++
		}
	}
	return total
}

func (dp DicePool) Values() []int {
	values := make([]int, dp.Count+1)
	for i := 0; i < len(values); i++ {
		values[i] = i
	}
	return values
}

func (dp DicePool) Hist() map[int]float64 {
	size := float64(dp.Dice.Size)
	gte := float64(dp.GTE)
	p := (size - gte + 1) / size
	h := histogram.Fixed(map[int]float64{0: 1.0 - p, 1: p})
	dd := make([]histogram.Histogram, dp.Count)
	for i := 0; i < dp.Count; i++ {
		dd[i] = h
	}
	hist := histogram.Aggregate(dd...).Hist()
	return hist
}
