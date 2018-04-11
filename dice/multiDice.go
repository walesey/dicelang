package dice

import "github.com/walesey/dicelang/histogram"

type MultiDice struct {
	Dice  Dice
	Count int
}

func (md MultiDice) Resolve() int {
	var total int
	for i := 0; i < md.Count; i++ {
		total += md.Dice.Resolve()
	}
	return total
}

func (md MultiDice) Values() []int {
	minValue := md.Count
	maxValue := md.Count * md.Dice.Size
	values := make([]int, maxValue-minValue+1)
	for i := 0; i < len(values); i++ {
		values[i] = i + minValue
	}
	return values
}

func (md MultiDice) Hist() map[int]float64 {
	dd := make([]histogram.Histogram, md.Count)
	for i := 0; i < md.Count; i++ {
		dd[i] = md.Dice
	}
	return histogram.Aggregate(dd...).Hist()
}
