package sorts

import "math/bits"

// Just some utility functions.

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

// find the smallest item greater than ref between left and right
func (sd *sortData) smallestGreaterThan(left, right, ref int) int {
	for left < right {
		m := median(left, right)
		if sd.lte(m, ref) {
			left = m + 1
		} else {
			right = m
		}
	}
	return left
}

// find the largest item less than ref between left and right
func (sd *sortData) largestLessThan(left, right, ref int) int {
	right++
	for left < right {
		m := median(left, right)
		if sd.less(m, ref) {
			left = m + 1
		} else {
			right = m
		}
	}
	return left
}

func (sd *sortData) rotateRightwards(left, right int) {
	for ; right > left; right-- {
		sd.Swap(right, right-1)
	}
}

func (sd *sortData) rotateLeftwards(left, right int) {
	for ; left < right; left++ {
		sd.Swap(left, left+1)
	}
}

func (sd *sortData) reverse(left, right int) {
	for left < right {
		sd.Swap(left, right)
		left++
		right--
	}
}

// swap the blocks at [a:a+count] and [b:b+count]
// Assumes ranges do not overlap.
func (sd *sortData) swapRange(a, b, count int) {
	for i := 0; i < count; i++ {
		sd.Swap(a+i, b+i)
	}
}

// swap internal blocks [left:mid] and [mid:right+1].
// Assumes that left < mid < right.
// If you are calling this with an offset if left+1 or right-1,
// stop and use rotateLeftwards or rotateRightwards instead.
func (sd *sortData) blockSwap(left, mid, right int) {
	i := mid - left
	j := right - mid + 1
	for i != j {
		if i > j {
			sd.swapRange(mid-i, mid, j)
			i -= j
		} else {
			sd.swapRange(mid-i, mid+j-i, i)
			j -= i
		}
	}
	sd.swapRange(mid-i, mid, i)
}
