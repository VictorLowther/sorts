package sorts

import "sort"

// Lifted from the standard library for benchmark and fooling around
// purposes, and adapted to use the sort primitives in this package.
//
// A shame that the way sort.Interface works more or less forces us to rely
// on an in-place merge algorithm -- powersort tends to collect runs that
// have unequal lengths, which doesn't work well with symMerge's preference
// for equal-ish length sides.  Perhaps a more complicated in-place merge
// sort would work better, but the ones based on block sort merging
// (http://itbe.hanyang.ac.kr/ak/papers/tamc2008.pdf, etc.) are much trickier
// to reason about and implement than symmerge is.
//
// Maybe generics will save us from the cruel tyranny forced in-place sorts.
//
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
func symMerge(sd sort.Interface, left, mid, oneAfterRight int) {
	// If either side only has one element to merge, just insert it
	// directly into place.
	var pivot int
	if mid-left == 1 {
		// one day, when the inliner is smarterer
		// pivot = largestLessThan(sd, mid, oneAfterRight-1, left)
		// rotateLeftwards(sd, left, pivot-1)
		pivot = mid
		for pivot < oneAfterRight {
			m := median(pivot, oneAfterRight)
			if sd.Less(m, left) {
				pivot = m + 1
			} else {
				oneAfterRight = m
			}
		}
		rotateLeftwards(sd, left, pivot-1)
		return
	}
	if oneAfterRight-mid == 1 {
		// Also when the inliner is smarter:
		// pivot = sd.smallestGreaterThan(left, mid, mid)
		// sd.rotateRightwards(pivot, mid)
		pivot = mid
		for left < pivot {
			m := median(left, pivot)
			if lte(sd, m, mid) {
				left = m + 1
			} else {
				pivot = m
			}
		}
		rotateRightwards(sd, left, mid)
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
		if !sd.Less(p-c, c) {
			start = c + 1
		} else {
			r = c
		}
	}

	end := n - start
	if start < mid && mid < end {
		blockSwap(sd, start, mid, end-1)
	}
	if left < start && start < pivot {
		symMerge(sd, left, start, pivot)
	}
	if pivot < end && end < oneAfterRight {
		symMerge(sd, pivot, end, oneAfterRight)
	}
}
