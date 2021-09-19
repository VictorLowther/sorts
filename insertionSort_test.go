package sorts

import (
	"math/rand"
	"testing"
)

func TestInsertionSorts(t *testing.T) {
	for _, sort := range []Sort{InsertionSort} {
		for i := 0; i < 100; i++ {
			vals := makeInts(rand.NewSource(rand.Int63()), 20, true)
			sort(vals)
			if failAt := vals.Sorted(); failAt < len(vals.v) {
				t.Errorf("iteration %d: Data not sorted at %d and %d (%0x,%04x) (%0x,%04x)", i, failAt-1, failAt,
					vals.v[failAt-1]>>16, vals.v[failAt-1]&0xffff, vals.v[failAt]>>16, vals.v[failAt]&0xffff)
				return
			}
		}
	}
}
