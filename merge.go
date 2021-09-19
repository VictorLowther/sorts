package sorts

import "sort"

// symMerge merges the two sorted subsequences data[a:m] and data[m:b] using
// the SymMerge algorithm from Pok-Son Kim and Arne Kutzner, "Stable Minimum
// Storage Merging by Symmetric Comparisons", in Susanne Albers and Tomasz
// Radzik, editors, Algorithms - ESA 2004, volume 3221 of Lecture Notes in
// Computer Science, pages 714-723. Springer, 2004.
//
// Let M = m-a and N = b-n. Wolog M < N.
// The recursion depth is bound by ceil(log(N+M)).
// The algorithm needs O(M*log(N/M + 1)) calls to data.Less.
// The algorithm needs O((M+N)*log(M)) calls to data.Swap.
//
// The paper gives O((M+N)*log(M)) as the number of assignments assuming a
// rotation algorithm which uses O(M+N+gcd(M+N)) assignments. The argumentation
// in the paper carries through for Swap operations, especially as the block
// swapping rotate uses only O(M+N) Swaps.
//
// symMerge assumes non-degenerate arguments: a < m && m < b.
// Having the caller check this condition eliminates many leaf recursion calls,
// which improves performance.
func (sd *sortData) symMerge(left, mid, oneAfterRight int) {
	var pivot int
	// If either side only has one element to merge, just insert it
	// directly into place.
	if mid-left == 1 {
		pivot = sd.largestLessThan(mid, oneAfterRight-1, left)
		sd.rotateLeftwards(left, pivot-1)
		return
	}
	if oneAfterRight-mid == 1 {
		pivot = sd.smallestGreaterThan(left, mid, mid)
		sd.rotateRightwards(pivot, mid)
		return
	}

	pivot = median(left, oneAfterRight)
	n := pivot + mid
	var start, r int
	if mid > pivot {
		start = n - oneAfterRight
		r = pivot
	} else {
		start = left
		r = mid
	}
	p := n - 1

	for start < r {
		c := int(uint(start+r) >> 1)
		if !sd.less(p-c, c) {
			start = c + 1
		} else {
			r = c
		}
	}

	end := n - start
	if start < mid && mid < end {
		sd.blockSwap(start, mid, end-1)
	}
	if left < start && start < pivot {
		sd.symMerge(left, start, pivot)
	}
	if pivot < end && end < oneAfterRight {
		sd.symMerge(pivot, end, oneAfterRight)
	}
}

func StdlibStable(vals []int, a Less) {
	sort.Stable(&sortData{data: vals, less: a})

}
