package sorts

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestPowerSort(t *testing.T) {
	for i := 0; i < 10000; i++ {
		vals := makeInts(rand.NewSource(7), i, true)
		Powersort(vals)
		if failAt := vals.Sorted(); failAt < len(vals.v) {
			t.Errorf("iteration %d: Data not sorted at %d and %d (%0x,%04x) (%0x,%04x)", i, failAt-1, failAt,
				vals.v[failAt-1]>>16, vals.v[failAt-1]&0xffff, vals.v[failAt]>>16, vals.v[failAt]&0xffff)
			return
		}
	}
	t.Logf("Data sorted")
}

type sortParams struct {
	totalLen int
	runs     bool
}

type sortMetric struct {
	name string
	s    Sort
}

func BenchmarkStableSorts(b *testing.B) {
	rand.Seed(int64(time.Now().UnixNano()))
	benchmarks := []sortParams{
		{10, false},
		{100, false},
		{100, true},
		{500, false},
		{500, true},
		{1000, false},
		{1000, true},
		{5000, false},
		{5000, true},
		{10000, false},
		{10000, true},
		{50000, false},
		{50000, true},
		{100000, true},
		{100000, false},
		{500000, false},
		{500000, true},
		//		{1000000, false},
		//		{1000000, true},
		//		{5000000, false},
		//		{5000000, true},
		//		{10000000, false},
		//		{10000000, true},
		//		{50000000, false},
		//		{50000000, true},
	}
	sorts := []*sortMetric{
		{name: "power", s: Powersort},
		{name: "stdlib", s: StdlibStable},
	}
	seed := rand.Int63()
	b.Logf("Using random seed %d", seed)
	src := rand.NewSource(seed)
	for _, bench := range benchmarks {
		for _, s := range sorts {
			b.Run(fmt.Sprintf("%s len %d, runs %v", s.name, bench.totalLen, bench.runs), func(bb *testing.B) {
				bb.StopTimer()
				vals := makeInts(src, bench.totalLen, bench.runs)
				bb.StartTimer()
				s.s(vals)
				bb.ReportMetric(float64(vals.cmps), "cmps")
				bb.ReportMetric(float64(vals.swaps), "swaps")
				bb.ReportMetric(float64(vals.maxRun), "maxrun")
			})
		}
	}
}
