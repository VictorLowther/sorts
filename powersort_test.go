package sorts

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func maskedLess(j []int) Less {
	return func(a, b int) bool {
		va, vb := j[a]>>16, j[b]>>16
		return va < vb
	}
}

func maskedSorted(j []int) int {
	if len(j) < 2 {
		return len(j)
	}
	less := maskedLess(j)
	for i := 1; i < len(j); i++ {
		if less(i, i-1) {
			return i
		}
		if !less(i-1, i) {
			if j[i]&0xffff >= j[i-1]&0xffff {
				return i
			}
		}
	}
	return len(j)
}

func makeInts(randSrc rand.Source, totalLen int, runs bool) []int {
	src := rand.New(randSrc)
	var pos int
	var mask uint16
	runLen := int16(1)
	if runs {
		runLen = int16(math.Sqrt(float64(totalLen))) + 1
	}
	res := []int{}
top:
	for {
		val := int(src.Int31()) << 32
		runCount := int16(src.Int31n(int32(runLen+1) + 1))
		// decide whether this is an ascending or descending run
		if src.Intn(2) == 0 {
			runCount = 0 - runCount
		}
		mask = 0
		for runCount != 0 {
			pos = len(res)
			if pos >= totalLen {
				break top
			}
			res = append(res, val)
			mask -= 1
			if runs {
				res[pos] += int(mask) + ((int(runCount) + src.Intn(2)) << 16)
			} else {
				res[pos] += int(mask) + (int(runCount) << 16)
			}
			if runCount > 0 {
				runCount--
			} else {
				runCount++
			}
		}
	}
	return res
}

func TestPowerSort(t *testing.T) {
	for i := 0; i < 10000; i++ {
		vals := makeInts(rand.NewSource(7), i, true)
		Powersort(vals, maskedLess(vals))
		if failAt := maskedSorted(vals); failAt < len(vals) {
			t.Errorf("iteration %d: Data not sorted at %d and %d (%0x,%04x) (%0x,%04x)", i, failAt-1, failAt,
				vals[failAt-1]>>16, vals[failAt-1]&0xffff, vals[failAt]>>16, vals[failAt]&0xffff)
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
		{50000, false},
		{100000, true},
		{100000, false},
		{500000, false},
		{500000, true},
		{1000000, false},
		{1000000, true},
		{5000000, false},
		{5000000, true},
		{10000000, false},
		{10000000, true},
		{50000000, false},
		{50000000, true},
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
			cmps := uint64(0)
			b.Run(fmt.Sprintf("%s len %d, runs %v", s.name, bench.totalLen, bench.runs), func(bb *testing.B) {
				bb.StopTimer()
				vals := makeInts(src, bench.totalLen, bench.runs)
				bb.StartTimer()
				less := maskedLess(vals)
				s.s(vals, func(a, b int) bool { cmps++; return less(a, b) })
				bb.ReportMetric(float64(cmps), "cmps")
			})
		}
	}
}
