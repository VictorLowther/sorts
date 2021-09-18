package sorts

import "testing"

func TestInsertionSorts(t *testing.T) {
	for _, sort := range []Sort{InsertionSort} {
		for i := 0; i < 100; i++ {
			vals := makeInts(10, 0)
			sort(vals, maskedLess(vals))
			if failAt := maskedSorted(vals); failAt < len(vals) {
				t.Errorf("Data not sorted at %d and %d (%08x) (%08x)", failAt-1, failAt, vals[failAt-1], vals[failAt])
				return
			}

		}
	}
}
