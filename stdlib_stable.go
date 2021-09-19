package sorts

import "sort"

// Also lifted from the standard library but adapted to
// use the sort primitives from this package. Used as a reference
// for benchmarking the efficiency of any other eventual stable
// sort implementations
func stdlibStable(data sort.Interface, n int) {
	blockSize := 20 // must be > 0
	a, b := 0, blockSize
	for b <= n {
		insertionSort(data, a, b-1, 0)
		a = b
		b += blockSize
	}
	insertionSort(data, a, n-1, 0)

	for blockSize < n {
		a, b = 0, 2*blockSize
		for b <= n {
			symMerge(data, a, a+blockSize, b)
			a = b
			b += 2 * blockSize
		}
		if m := a + blockSize; m < n {
			symMerge(data, a, m, n)
		}
		blockSize *= 2
	}
}

func StdlibStable(v sort.Interface) { stdlibStable(v, v.Len()) }
