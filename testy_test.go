package sorts

import (
	"math"
	"math/rand"
)

// maskedInts is used for unit tests and benchmarks.
// The idea is to generate a mostly-random slice of integers
// with some ascending and descending runs, with short stretches
// of equal values.  The least significant 16 bits are reserved for
// an always-descending sequence that we can use to detect any instability
// in a supposedly stable sort.
type maskedInts struct {
	v      []int
	swaps  uint64
	cmps   uint64
	maxRun int
}

func (m *maskedInts) Len() int {
	return len(m.v)
}

func (m *maskedInts) Less(a, b int) bool {
	m.cmps++
	return m.v[a]>>16 < m.v[b]>>16
}

func (m *maskedInts) Swap(a, b int) {
	m.swaps++
	m.v[a], m.v[b] = m.v[b], m.v[a]
}

func (m *maskedInts) Sorted() int {
	if len(m.v) < 2 {
		return len(m.v)
	}
	for i := 1; i < len(m.v); i++ {
		if m.Less(i, i-1) {
			return i
		}
		if !m.Less(i-1, i) {
			if m.v[i]&0xffff >= m.v[i-1]&0xffff {
				return i
			}
		}
	}
	return len(m.v)
}

func makeInts(randSrc rand.Source, totalLen int, runs bool) *maskedInts {
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
	return &maskedInts{v: res, maxRun: int(runLen)}
}
