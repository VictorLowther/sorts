package sorts

func (sd *sortData) binarySearchBackwards(left, right, ref int) int {
	mid := 0
	for left < right {
		mid = int(uint(left+right)>>1) + 1
		if sd.gt(ref, mid) {
			right = mid - 1
		} else {
			left = mid
		}
	}
	return right
}

func (sd *sortData) binarySearchForwards(left, right, ref int) int {
	mid := 0
	for left < right {
		mid = int(uint(left+right) >> 1)
		if sd.less(ref, mid) {
			right = mid
		} else {
			left = mid + 1
		}

	}
	return left
}

func (sd *sortData) reverse(left, right int) {
	for left < right {
		sd.data[left], sd.data[right] = sd.data[right], sd.data[left]
		left++
		right--
	}
}

func (sd *sortData) insertionSort(first, last, sortedAt int) {
	if sortedAt < first {
		sortedAt = first
	}
	sortedAt++
	for ; sortedAt <= last; sortedAt++ {
		for j := sortedAt; j > first && sd.less(j, j-1); j-- {
			sd.data[j], sd.data[j-1] = sd.data[j-1], sd.data[j]
		}
	}
}

// InsertionSort is a simple insertion sort.  It is very cache-friendly
// for small amounts of data, but the runtime costs are terrible for more
// than 16 or so items
func InsertionSort(vals []int, a Less) {
	(&sortData{data: vals, less: a}).insertionSort(0, len(vals)-1, 0)
}
