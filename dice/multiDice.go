package dice

import "math"

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
	nbPossibilities := math.Pow(float64(md.Dice.Size), float64(md.Count))
	// recursively add probablities of every possible result
	hist := make(map[int]float64)
	var recur func(total, digit int)
	recur = func(total, digit int) {
		for _, v := range md.Dice.Values() {
			if digit >= md.Count {
				k := total + v
				hist[k] = hist[k] + 1
			} else {
				recur(total+v, digit+1)
			}
		}
	}
	recur(0, 1)
	for k, v := range hist {
		hist[k] = v / nbPossibilities
	}
	return hist
}
