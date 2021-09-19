// sorts is a playground for experimenting with various stable sort algorithms.
package sorts

import "sort"

// Less is any function that can compare to see of the
// array element at the first position is less than the
type Less func(int, int) bool

// Sort is any function that takes an array of integers and
// a Less function and sorts that array based on it.
type Sort func([]int, Less)

type internalSort func(*sortData, int, int, int)

// Interface uses fn to sort a collection that conforms
// to the standard library sort.Interface.
//
// It allocates a slice containing the indexes of the
// data to be sorted, sorts that slice based on the
// s.Less, and then swaps the underlying data into
// its final location.
func Interface(fn Sort, s sort.Interface) {
	indexes := make([]int, s.Len())
	for i := 0; i < len(indexes); i++ {
		indexes[i] = i
	}

	fn(indexes, func(i, j int) bool {
		return s.Less(i, j)
	})

	for i := 0; i < len(indexes); i++ {
		j := indexes[i]
		if j == 0 {
			continue
		}
		for k := i; j != i; {
			s.Swap(j, k)
			k, j, indexes[j] = j, indexes[j], 0
		}
	}
}

// All other sorts should build off of the methods
// this struct provides.
type sortData struct {
	less Less
	data []int
}

func (s *sortData) gte(a, b int) bool {
	return !s.less(a, b)
}

func (s *sortData) gt(a, b int) bool {
	return s.less(b, a)
}

func (s *sortData) lte(a, b int) bool {
	return !s.less(b, a)
}

func (s *sortData) Len() int {
	return len(s.data)
}

func (s *sortData) Less(a, b int) bool {
	return s.less(a, b)
}

func (s *sortData) Swap(a, b int) {
	s.data[a], s.data[b] = s.data[b], s.data[a]
}
