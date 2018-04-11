package dice

import "math/rand"

type Dice struct {
	Size    int
	Weights []float64
}

func (d Dice) Resolve() int {
	return rand.Intn(d.Size) + 1
}

func (d Dice) Values() []int {
	values := make([]int, d.Size)
	for i := 0; i < len(values); i++ {
		values[i] = i + 1
	}
	return values
}

func (d Dice) Hist() map[int]float64 {
	hist := make(map[int]float64)
	for i, v := range d.Values() {
		if d.Weights != nil && len(d.Weights) == d.Size {
			hist[v] = d.Weights[i]
		} else {
			hist[v] = 1.0 / float64(d.Size)
		}
	}
	return hist
}
