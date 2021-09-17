package sorts

import (
	"math"
	"sort"
)

// Powersort, an adaptive mergesort, described here: https://arxiv.org/pdf/1805.04154.pdf
// The key improvement powersort makes over the standard library merge sort is in carefully
// choosing when to merge adjacent runs of already sorted data based on their calculated
// position in a nearly-optimal binary merge tree.

var (
	minRunLen = 16
)

// reverse all items in a range. Used when we detect a run
// of descending data.
func reverse(a sort.Interface, start, end int) {
	for start < end {
		a.Swap(start, end)
		start++
		end--
	}
}

// since we don't already have it when we need it...
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
func findOrCreateRun(a sort.Interface, start, last int) (sortedTo int) {
	sortedTo = start
	if sortedTo == last {
		return
	}
	sortedTo++
	if a.Less(sortedTo, sortedTo-1) {
		// Assume we are in a descending run.
		// Stop at the first item that is not in strict descending order.
		// Reverse the list that we found that was in strict descending order.
		for ; sortedTo < last; sortedTo++ {
			if !a.Less(sortedTo+1, sortedTo) {
				break
			}
		}
		reverse(a, start, sortedTo)
	} else {
		// Assume we are in a weakly ascending run.
		// Stop at the first item that is not in weakly ascensing order.
		for ; sortedTo < last; sortedTo++ {
			if a.Less(sortedTo+1, sortedTo) {
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
		insertionSort(a, start, last, sortedTo)
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
func nodePower(totalLen, startA, startB, endB int) int {
	relativePowerA := startA + startB
	relativePowerB := startB + endB + 1
	commonBits := 0
	var digitA, digitB bool
	for {
		digitA, digitB = relativePowerA >= totalLen, relativePowerB >= totalLen
		if digitA != digitB {
			break
		}
		commonBits++
		if digitA {
			relativePowerA -= totalLen
		}
		if digitB {
			relativePowerB -= totalLen
		}
		relativePowerA <<= 1
		relativePowerB <<= 1

	}
	return commonBits
}

// All this stuff will get inlined, since they are all leaf calls.
type runStack [][2]int

func (r runStack) entryAt(i int) bool {
	return !(r[i][0] == 0 && r[i][1] == 0)
}

func (r runStack) setAt(i, first, last int) {
	r[i][0], r[i][1] = first, last
}

func (r runStack) mergeAt(a sort.Interface, i, last int) int {
	symMerge(a, r[i][0], r[i][1]+1, last)
	return r[i][0]
}

func PowerSort(a sort.Interface) {
	var top, firstA, lastA, firstB, lastB, power, i, totalLen int
	totalLen = a.Len()
	// Entries in the runs are indexed by their node power.
	// Node power will never be more than log2(totalLen) + 1
	// Therefore, the runstack needs to contain (log2(totalLen)+1) entries
	runs := make(runStack, int(math.Log2(float64(totalLen)))+1)
	lastA = findOrCreateRun(a, 0, totalLen-1)
	for lastA < totalLen-1 {
		firstB, lastB = lastA+1, findOrCreateRun(a, lastA+1, totalLen-1)
		power = nodePower(totalLen, firstA, firstB, lastB)
		for i = top; i > power; i-- {
			if runs.entryAt(i) {
				firstA = runs.mergeAt(a, i, lastA+1)
				runs.setAt(i, 0, 0)
			}
		}
		runs.setAt(power, firstA, lastA)
		top = power
		firstA, lastA = firstB, lastB
	}
	for i = top; i >= 0; i-- {
		if runs.entryAt(i) {
			runs.mergeAt(a, i, totalLen)
		}
	}
}
