package histogram

import (
	"sort"

	"github.com/walesey/dicelang/util"
)

// HistogramColumn - stores a Value/Probablity
type HistogramColumn struct {
	V int
	P float64
}

type ByValue []HistogramColumn

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValue) Less(i, j int) bool { return a[i].V < a[j].V }

func FormatHistogram(hist Histogram) []HistogramColumn {
	h := hist.Hist()
	hList := make([]HistogramColumn, 0, len(h))
	for v, p := range h {
		hList = append(hList, HistogramColumn{V: v, P: p})
	}
	sort.Sort(ByValue(hList))
	return hList
}

func RoundHistogram(hist Histogram, places int) Histogram {
	h := hist.Hist()
	for v, p := range h {
		newP := util.Round(p, .5, places)
		if newP == 0 {
			delete(h, v)
		} else {
			h[v] = newP
		}
	}
	return Fixed(h)
}
