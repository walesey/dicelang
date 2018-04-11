package dice

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
	p := float64(dp.Dice.Size-dp.GTE+1) / float64(dp.Dice.Size)
	weightedD2 := Dice{Size: 2, Weights: []float64{(1.0 - p), p}}
	md := MultiDice{Dice: weightedD2, Count: dp.Count}
	h := md.Hist()
	hist := make(map[int]float64)
	for k, v := range h {
		hist[k-dp.Count] = v
	}
	return hist
}
