package sorts
import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

type justInts []uint

func (j justInts) Len() int {
	return len(j)
}

func (j justInts) Less(a, b int) bool {
	return j[a] < j[b]
}

func (j justInts) Swap(a, b int) {
	j[a], j[b] = j[b], j[a]
}

func (j justInts) Sorted() int {
	if len(j) < 2 {
		return len(j)
	}
	for i := 1; i < len(j); i++ {
		if j.Less(i, i-1) {
			return i
		}
	}
	return len(j)
}

type maskedInts []uint

func (j maskedInts) Len() int {
	return len(j)
}

func (j maskedInts) Less(a, b int) bool {
	va, vb := j[a]>>8, j[b]>>8
	return va < vb
}

func (j maskedInts) Swap(a, b int) {
	j[a], j[b] = j[b], j[a]
}

func (j maskedInts) Sorted() int {
	if len(j) < 2 {
		return len(j)
	}
	for i := 1; i < len(j); i++ {
		if j.Less(i, i-1) {
			return i
		}
		if !j.Less(i-1, i) {
			if j[i]&0xff >= j[i-1]&0xff {
				return i
			}
		}
	}
	return len(j)
}

func makeInts(totalLen, runLen, maskLen int) []uint {
	var pos int
	res := []uint{}
top:
	for {
		val := uint(rand.Int()) << 16
		runCount := rand.Intn(runLen+1) + 1
		if rand.Intn(2) == 0 {
			runCount = 0 - runCount
		}
		for runCount != 0 {
			for maskCount := rand.Intn(maskLen+1) + 1; maskCount > 0; maskCount-- {
				pos = len(res)
				if pos >= totalLen {
					break top
				}
				res = append(res, val)
				res[pos] |= uint(maskCount&0xff + ((runCount & 0xff) << 8))
			}
			if runCount < 0 {
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
	vals := maskedInts(makeInts(100, 10, 10))
	PowerSort(vals)
	if failAt := vals.Sorted(); failAt < len(vals) {
		t.Errorf("Data not sorted at %d and %d (%08x) (%08x)", failAt-1, failAt, vals[failAt-1], vals[failAt])
		return
	}
	t.Logf("Data sorted")
}

type sortParams struct {
	totalLen int
	runLen   int
	maskLen  int
}

func BenchmarkStableSorts(b *testing.B) {
	rand.Seed(int64(time.Now().UnixNano()))
	benchmarks := []sortParams{
		{10, 0, 0},
		{100, 0, 0},
		{100, 10, 5},
		{500, 0, 0},
		{500, 30, 10},
		{1000, 0, 0},
		{1000, 10, 5},
		{5000, 0, 0},
		{5000, 10, 10},
		{5000, 50, 10},
		{10000, 0, 0},
		{10000, 10, 5},
		{50000, 0, 0},
		{50000, 10, 10},
		{50000, 50, 10},
		{100000, 0, 0},
		{100000, 10, 5},
		{500000, 0, 0},
		{500000, 10, 10},
		{500000, 50, 10},
	}
	for _, bench := range benchmarks {
		vals := makeInts(bench.totalLen, bench.runLen, bench.maskLen)
		for _, minRunLen = range []int{16} {
			b.Run(fmt.Sprintf("power minrun %d len %d, runlen %d, masklen %d", minRunLen, bench.totalLen, bench.runLen, bench.maskLen), func(bb *testing.B) {
				PowerSort(maskedInts(append([]uint{}, vals...)))
			})
		}
		b.Run(fmt.Sprintf("stable len %d, runlen %d, masklen %d", bench.totalLen, bench.runLen, bench.maskLen), func(bb *testing.B) {
			sort.Stable(maskedInts(append([]uint{}, vals...)))
		})
	}

}
