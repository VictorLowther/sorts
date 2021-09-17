package sorts

import "sort"

func binarySearchForwards(a sort.Interface, left, right, ref int) int {
	mid := 0
	for left < right {
		mid = int(uint(left+right)>>1) + 1
		if a.Less(mid, ref) {
			left = mid
		} else {
			right = mid - 1
		}
	}
	return right
}

func moveLeftmostElementForwards(a sort.Interface, left, right int) {
	target := binarySearchForwards(a,left,right,left)
	for ; left < target; left++ {
		a.Swap(left, left+1)
	}
}

func binarySearchBackwards(a sort.Interface, left, right, ref int) int {
	mid := 0
	for left < right {
		mid = int(uint(left+right) >> 1)
		if a.Less(ref, mid) {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
}

func moveRightmostElementBackwards(a sort.Interface, left, right int) {
	target := binarySearchBackwards(a,left,right,right)
	for ; right > target; right-- {
		a.Swap(right, right-1)
	}
}

// insertionSort is a binary insertion sort, geared to sorting part of a larger list.
// first and last are the first and last indexes in the data to sort in,
// and start is the first
// element between lo and hi that we already know to be sorted.
func insertionSort(a sort.Interface, first, last, sortedAt int) {
	if sortedAt < first {
		sortedAt = first
	}
	sortedAt++
	for ; sortedAt <= last; sortedAt++ {
		moveRightmostElementBackwards(a, first, sortedAt)
	}
}

func InsertionSort(a sort.Interface) {
	insertionSort(a,0,a.Len(),0)
}
