package sorts

import "sort"

// Bog-standard insertion sort, with an added bonus feature
// to indicate the rightmost sorted position.
// powersort uses that feature to leverage existing runs.
func insertionSort(sd sort.Interface, left, right, sortedAt int) {
	if sortedAt < left {
		sortedAt = left
	}
	sortedAt++
	for ; sortedAt <= right; sortedAt++ {
		for j := sortedAt; j > left && sd.Less(j,j-1); j-- {
			sd.Swap(j,j-1)
		}
	}
}

// InsertionSort is a simple insertion sort.  It is very cache-friendly
// for small amounts of data, but the runtime costs are terrible for more
// than 32 or so items
func InsertionSort(v sort.Interface) {
	insertionSort(v,0,v.Len()-1,0)
}
