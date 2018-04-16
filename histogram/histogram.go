package histogram

import (
	"math/rand"
)

type Histogram interface {
	Resolve() int
	Hist() map[int]float64
}

type Fixed map[int]float64

func (fh Fixed) Resolve() int {
	hist := fh.Hist()
	rnb := rand.Float64()
	var f float64
	for val, prob := range hist {
		if rnb > f && rnb <= f+prob {
			return val
		}
		f += prob
	}
	return -1
}

func (fh Fixed) Hist() map[int]float64 {
	return fh
}

// Aggregate [4d6.4+, 4d6.6+]
func Aggregate(histograms ...Histogram) Histogram {
	if len(histograms) == 0 {
		return nil
	}
	hist := histograms[0]
	for i := 1; i < len(histograms); i++ {
		hist = aggregate(hist, histograms[i])
	}
	return hist
}

func aggregate(h1, h2 Histogram) Histogram {
	hist1 := h1.Hist()
	hist2 := h2.Hist()
	hist := make(map[int]float64)
	for v1, p1 := range hist1 {
		for v2, p2 := range hist2 {
			v := v1 + v2
			p := p1 * p2
			if current, ok := hist[v]; ok {
				hist[v] = current + p
			} else {
				hist[v] = p
			}
		}
	}
	return Fixed(hist)
}

// Multiply 2d3.add.d6.4+
func Multiply(histograms ...Histogram) Histogram {
	if len(histograms) == 0 {
		return nil
	}
	hist := histograms[0]
	for i := 1; i < len(histograms); i++ {
		hist = multiply(hist, histograms[i])
	}
	return hist
}

func multiply(h1, h2 Histogram) Histogram {
	formattedH1 := FormatHistogram(h1)
	hist := make(map[int]float64)
	h2Aggregate := h2
	count := 1
	for _, column := range formattedH1 {
		v1, p1 := column.V, column.P
		if v1 == 0 {
			continue
		}
		for i := count; i < v1; i++ {
			h2Aggregate = aggregate(h2Aggregate, h2)
			count++
		}
		for v2, p2 := range h2Aggregate.Hist() {
			v := v2
			p := p1 * p2
			if current, ok := hist[v]; ok {
				hist[v] = current + p
			} else {
				hist[v] = p
			}
		}
	}
	return Fixed(hist)
}
