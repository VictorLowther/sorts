package sorts

func (sd *sortData) insertionSort(left, right, sortedAt int) {
	if sortedAt < left {
		sortedAt = left
	}
	sortedAt++
	for ; sortedAt <= right; sortedAt++ {
		target := sd.smallestGreaterThan(left, sortedAt, sortedAt)
		sd.rotateRightwards(target, sortedAt)
	}
}

// InsertionSort is a simple insertion sort.  It is very cache-friendly
// for small amounts of data, but the runtime costs are terrible for more
// than 32 or so items
func InsertionSort(vals []int, a Less) {
	(&sortData{data: vals, less: a}).insertionSort(0, len(vals)-1, 0)
}
