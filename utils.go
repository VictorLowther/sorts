package sorts

import (
	"math/bits"
	"sort"
)

// Just some useful utility functions. Most of them
// wind up being inlined.
func lte(s sort.Interface, a, b int) bool {
	return !s.Less(b, a)
}

func log2(s uint64) int {
	return 64 - bits.LeadingZeros64(s) - 1
}

func median(a, b int) int {
	return int(uint(a+b) >> 1)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func rotateRightwards(sd sort.Interface, left, right int) {
	for ; right > left; right-- {
		sd.Swap(right, right-1)
	}
}

func rotateLeftwards(sd sort.Interface, left, right int) {
	for ; left < right; left++ {
		sd.Swap(left, left+1)
	}
}

func reverse(sd sort.Interface, left, right int) {
	for left < right {
		sd.Swap(left, right)
		left++
		right--
	}
}

// swap the blocks at [a:a+count] and [b:b+count]
// Assumes ranges do not overlap.
func swapRange(sd sort.Interface, a, b, count int) {
	for i := 0; i < count; i++ {
		sd.Swap(a+i, b+i)
	}
}

// It would be great if the compiler were able to inline these.
// Just saying.

// swap internal blocks [left:mid] and [mid:right+1] using
// the smallest number of swap operations.  Assumes
// that left < mid < right.
func blockSwap(sd sort.Interface, left, mid, right int) {
	i := mid - left
	j := right - mid + 1
	for i != j {
		if i > j {
			swapRange(sd, mid-i, mid, j)
			i -= j
		} else {
			swapRange(sd, mid-i, mid+j-i, i)
			j -= i
		}
	}
	swapRange(sd, mid-i, mid, i)
}

// find the smallest item greater than ref between left and right
func smallestGreaterThan(sd sort.Interface, left, right, ref int) int {
	for left < right {
		m := median(left, right)
		if lte(sd, m, ref) {
			left = m + 1
		} else {
			right = m
		}
	}
	return left
}

// find the largest item less than ref between left and right
func largestLessThan(sd sort.Interface, left, right, ref int) int {
	right++
	for left < right {
		m := median(left, right)
		if sd.Less(m, ref) {
			left = m + 1
		} else {
			right = m
		}
	}
	return left
}
