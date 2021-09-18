package sorts

import (
	"math"
)

// Powersort, an adaptive mergesort, described here: https://arxiv.org/pdf/1805.04154.pdf
// The key improvement powersort makes over the standard library merge sort is in carefully
// choosing when to merge adjacent runs of already sorted data based on their calculated
// position in a nearly-optimal binary merge tree.

var (
	minRunLen = 32
)

type powerSort struct {
	sortData
	mergeBuf []int
}

// Find or create an ascending run in the given range.
// This will also handle reversing descending runs.
//
// start is the index of the first item to consider, end is the
// index of the last item to consider.
//
// sortedTo is the last index that is part of the run.
//
// If we found a natural run that is less than minRunLen, we will
// extend it to minRunLen using an insertion sort .
func (p *powerSort) findOrCreateRun(start, last int) (sortedTo int) {
	sortedTo = start
	if sortedTo == last {
		return
	}
	sortedTo++
	if p.less(sortedTo, sortedTo-1) {
		// Assume we are in a descending run.
		// Stop at the first item that is not in strict descending order.
		// Reverse the list that we found that was in strict descending order.
		for ; sortedTo < last; sortedTo++ {
			if !p.less(sortedTo+1, sortedTo) {
				break
			}
		}
		p.reverse(start, sortedTo)
	} else {
		// Assume we are in a weakly ascending run.
		// Stop at the first item that is not in weakly ascensing order.
		for ; sortedTo < last; sortedTo++ {
			if p.less(sortedTo+1, sortedTo) {
				break
			}
		}
	}
	if start+minRunLen-1 <= last &&
		sortedTo < last &&
		sortedTo-start < minRunLen-1 {
		// Found a too-short natural run that we can extend.
		// Extend and sort it using an insertion sort.
		last = min(last, start+minRunLen-1)
		p.insertionSort(start, last, sortedTo)
		sortedTo = last
	}
	return
}

// This is straight up black magic voodoo that calculates the hypothetical
// depth in a nearly optimal binary merge tree this run would have.
// The idea is to defer merges until we are ascending in the tree, and at that
// point the merges that are queued up in the runs list are in the more or less
// optimal order to be merged.
//
// Adapted from nodePowerBitwise in
// https://github.com/sebawild/nearly-optimal-mergesort-code/blob/master/src/wildinter/net/mergesort/PowerSort.java
func (p *powerSort) nodePower(startA, startB, endB int) int {
	relativePowerA := startA + startB
	relativePowerB := startB + endB + 1
	commonBits := 0
	var digitA, digitB bool
	for {
		digitA, digitB = relativePowerA >= len(p.data), relativePowerB >= len(p.data)
		if digitA != digitB {
			break
		}
		commonBits++
		if digitA {
			relativePowerA -= len(p.data)
		}
		if digitB {
			relativePowerB -= len(p.data)
		}
		relativePowerA <<= 1
		relativePowerB <<= 1

	}
	return commonBits
}

func (p *powerSort) mergeAt(left, mid, right int) {
	if !p.less(mid+1, mid) {
		// Already merged
		return
	}
	left = p.binarySearchForwards(left, mid-1, mid)
	right = p.binarySearchBackwards(mid, right, mid-1)
	start := left
	p.mergeBuf = p.mergeBuf[:0]
	mid++
	pivot := mid
	for left < pivot && mid < right+1 {
		if p.less(mid, left) {
			p.mergeBuf = append(p.mergeBuf, p.data[mid])
			mid++
		} else {
			p.mergeBuf = append(p.mergeBuf, p.data[left])
			left++
		}
	}
	if left < pivot {
		p.mergeBuf = append(p.mergeBuf, p.data[left:pivot]...)
	} else {
		p.mergeBuf = append(p.mergeBuf, p.data[mid:right+1]...)
	}
	copy(p.data[start:], p.mergeBuf)

}

func (p *powerSort) sort() {
	var top, firstA, lastA, firstB, lastB, power, i, totalLen int
	totalLen = len(p.data)
	// Entries in the runs are indexed by their node power.
	// Node power will never be more than log2(totalLen) + 1
	// Therefore, the runstack needs to contain (log2(totalLen)+1) entries
	runs := make([][2]int, int(math.Log2(float64(totalLen)))+1)
	entryAt := func(i int) bool { return !(runs[i][0] == 0 && runs[i][1] == 0) }
	lastA = p.findOrCreateRun(0, totalLen-1)
	for lastA < totalLen-1 {
		firstB, lastB = lastA+1, p.findOrCreateRun(lastA+1, totalLen-1)
		power = p.nodePower(firstA, firstB, lastB)
		for i = top; i > power; i-- {
			if entryAt(i) {
				p.mergeAt(runs[i][0], runs[i][1], lastA)
				firstA = runs[i][0]
				runs[i][0], runs[i][1] = 0, 0
			}
		}
		runs[power][0], runs[power][1] = firstA, lastA
		top = power
		firstA, lastA = firstB, lastB
	}
	for i = top; i >= 0; i-- {
		if entryAt(i) {
			p.mergeAt(runs[i][0], runs[i][1], totalLen-1)
		}
	}
}

func Powersort(vals []int, a Less) {
	ps := &powerSort{}
	ps.data = vals
	ps.less = a
	ps.sort()
}
