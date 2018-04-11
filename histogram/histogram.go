package histogram

import (
	"math"
)

type Histogram interface {
	Resolve() int
	Hist() map[int]float64
}

type Fixed map[int]float64

func (fh Fixed) Resolve() int {
	return 0 //TODO
}

func (fh Fixed) Hist() map[int]float64 {
	return fh
}

// Aggregate [4d6.4+, 4d6.6+]
func Aggregate(histograms ...Histogram) Histogram {
	hist := make(map[int]float64)
	var recur func(prob float64, total, i int)
	recur = func(prob float64, total, i int) {
		for v, p := range histograms[i].Hist() {
			if i >= len(histograms)-1 {
				k := total + v
				if current, ok := hist[k]; ok {
					hist[k] = current + prob*p
				} else {
					hist[k] = prob * p
				}
			} else {
				recur(prob*p, total+v, i+1)
			}
		}
	}
	recur(1, 0, 0)
	return Fixed(hist)
}

// Multiply 2d3.add.d6.4+
func Multiply(histograms ...Histogram) Histogram {
	hist := make(map[int]float64)
	var recur func(prob float64, total, i int)
	recur = func(prob float64, total, i int) {
		hList := make([]Histogram, total)
		for j := 0; j < total; j++ {
			hList[j] = histograms[i]
		}
		histogram := Aggregate(hList...)
		for v, p := range histogram.Hist() {
			if i >= len(histograms)-1 {
				if current, ok := hist[v]; ok {
					hist[v] = current + prob*p
				} else {
					hist[v] = prob * p
				}
			} else {
				recur(prob*p, v, i+1)
			}
		}
	}
	recur(1, 1, 0)
	return Fixed(hist)
}

func Invert(histogram Histogram) Histogram {
	hist := make(map[int]float64)
	for k, v := range histogram.Hist() {
		hist[k] = 1.0 - v
	}
	return Fixed(hist)
}

func RoundHistogram(h map[int]float64) map[int]float64 {
	unit := 0.001
	for k, v := range h {
		h[k] = math.Round(v/unit) * unit
	}
	return h
}
