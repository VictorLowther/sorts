package sorts

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func maskedLess(j []int) Less {
	return func(a, b int) bool {
		va, vb := j[a]>>8, j[b]>>8
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
			if j[i]&0xff >= j[i-1]&0xff {
				return i
			}
		}
	}
	return len(j)
}

func makeInts(totalLen int, runLen int16) []int {
	var pos int
	var mask byte
	res := []int{}
top:
	for {
		val := int(rand.Int31()) << 32
		runCount := int16(rand.Int31n(int32(runLen+1) + 1))
		if rand.Intn(2) == 0 {
			runCount = 0 - runCount
		}
		for runCount != 0 {
			pos = len(res)
			if pos >= totalLen {
				break top
			}
			res = append(res, val)
			mask -= 1
			res[pos] += int(mask) + (int(runCount) << 8)
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
	rand.Seed(int64(time.Now().UnixNano()))
	for i := 1; i < 100; i++ {
		vals := makeInts(100, 10)
		Powersort(vals, maskedLess(vals))
		if failAt := maskedSorted(vals); failAt < len(vals) {
			t.Errorf("Data not sorted at %d and %d (%08x) (%08x)", failAt-1, failAt, vals[failAt-1], vals[failAt])
			return
		}
	}
	t.Logf("Data sorted")
}

type sortParams struct {
	totalLen int
	runLen   int16
}

func BenchmarkStableSorts(b *testing.B) {
	rand.Seed(int64(time.Now().UnixNano()))
	benchmarks := []sortParams{
		{10, 0},
		{100, 0},
		{100, 10},
		{500, 0},
		{500, 30},
		{1000, 0},
		{1000, 10},
		{5000, 0},
		{5000, 10},
		{5000, 50},
		{10000, 0},
		{10000, 10},
		{50000, 0},
		{50000, 10},
		{50000, 50},
		{100000, 0},
		{100000, 10},
		{500000, 0},
		{500000, 10},
		{500000, 50},
		{1000000, 0},
		{1000000, 100},
		{5000000, 0},
		{5000000, 100},
		{5000000, 200},
		{10000000, 0},
		{100000, 200},
		{500000, 0},
		{50000000, 200},
		{50000000, 200},
	}
	for _, bench := range benchmarks {
		b.Run(fmt.Sprintf("power minrun %d len %d, runlen %d", minRunLen, bench.totalLen, bench.runLen), func(bb *testing.B) {
			bb.StopTimer()
			vals := makeInts(bench.totalLen, bench.runLen)
			bb.StartTimer()
			Powersort(vals, maskedLess(vals))
		})
		b.Run(fmt.Sprintf("stable len %d, runlen %d", bench.totalLen, bench.runLen), func(bb *testing.B) {
			bb.StopTimer()
			vals := makeInts(bench.totalLen, bench.runLen)
			bb.StartTimer()
			StdlibStable(vals, maskedLess(vals))
		})
	}

}
