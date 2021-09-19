package sorts

import "sort"

var (
	// Chosen to work well-ish enough. Too small, and symMerge
	// moves too much data, and too big, the On^2 nature of
	// insertion sort starts eating you alive.
	minRunLen = 16
)

// Find or create an ascending run in the given range.
// This will also handle reversing descending runs.
//
// left is the index of the first item to consider, right is the
// index of the last item to consider.
//
// last is the last index that is part of the run.
//
// If we found a natural run that is less than minRunLen, we will
// extend it to minRunLen using an insertion sort .
func findOrCreateRun(p sort.Interface, left, right int) (last int) {
	last = left
	if last == right {
		return
	}
	last++
	if p.Less(last, last-1) {
		// Assume we are in a descending run.
		// Stop at the first item that is not in strict descending order.
		// Reverse the list that we found that was in strict descending order.
		for ; last < right; last++ {
			if !p.Less(last+1, last) {
				break
			}
		}
		reverse(p, left, last)
	} else {
		// Assume we are in a weakly ascending run.
		// Stop at the first item that is not in weakly ascensing order.
		for ; last < right; last++ {
			if p.Less(last+1, last) {
				break
			}
		}
	}

	if last-left < minRunLen-1 {
		// Found a too-short natural run that we can extend.
		// Extend and sort it using an insertion sort.
		right = min(right, left+minRunLen-1)
		insertionSort(p, left, right, last)
		last = right
	}
	return
}

// nodePower is adapted from nodePowerBitwise in
// https://github.com/sebawild/nearly-optimal-mergesort-code/blob/master/src/wildinter/net/mergesort/PowerSort.java
//
// I have not spent enough time wrapping my head around what nodePower does
// to determine where two adjacent runs are on-the-fly in a nearly optimal
// binary search tree built out of the runs we want to search, but Tim Peters
// did when he modified TimSort in Python to use powersort's merge strategy
// in https://github.com/python/cpython/commit/5cb4c672d855033592f0e05162f887def236c00a
//
// Since his explanation is so much more readable than what is in the paper,
// I have copied it nearly verbatim below, with changes for the different
// function and variable names, and platform limitations:
//
// (from Tim Peters)
//
// The nodePower function computes a run's "power". Say two adjacent runs
// begin at index startA. The first runs from left to mid-1, and the second
// runs from mid to right.  The lengths of the first run is n1 := mid-left,
// and the length of the second run is n2 := right - mid + 1.
// The list has total length n.
//
// The "power" of the first run is a small integer, the depth of the node
// connecting the two runs in an ideal binary merge tree, where power 1 is the
// root node, and the power increases by 1 for each level deeper in the tree.
//
// The power is the least integer L such that the "midpoint interval" contains
// a rational number of the form J/2**L. The midpoint interval is the semi-
// closed interval:
//
//     ((left + n1/2)/n, (mid + n2/2)/n]
//
// Yes, that's brain-busting at first ;-) Concretely, if (left + n1/2)/n and
// (mid + n2/2)/n are computed to infinite precision in binary, the power L is
// the first position at which the 2**-L bit differs between the expansions.
// Since the left end of the interval is less than the right end, the first
// differing bit must be a 0 bit in the left quotient and a 1 bit in the right
// quotient.
//
// nodePower emulates these divisions, 1 bit at a time, using comparisons,
// subtractions, and shifts in a loop.
//
// You'll notice the paper uses an O(1) method instead that relies on
// integer divison on an integer type twice as wide as needed to hold the
// list length, and a fast method of counting the number of leading
// zeros in an integer.  Go has the latter in bits.LeadingZeros, but on
// 64 bit platforms Go's native integer type is already an int64, and we do not
// have a handy int128 lying around.
//
// But since runs in our algorithm are almost never very short, the once-per-run
// overhead of nodePower seems lost in the noise.
//
// (end from Tim Peters)
//
// For the interested, the Java code that does this with no loops is as follows:
//
//    private static int nodePower(int size, int left, int mid, int right) {
//	    int twoN = size << 1; // 2*n
//	    long l = left + mid;
//	    long r = mid + right + 1;
//	    int a = (int) ((l << 31) / twoN);
//	    int b = (int) ((r << 31) / twoN);
//	    return Integer.numberOfLeadingZeros(a ^ b);
//    }
//
func nodePower(length, left, mid, right int) int {
	leftPower := uint(left) + uint(mid)
	rightPower := uint(mid) + uint(right) + 1
	size := uint(length)
	power := 0
	var digitA, digitB bool
	for {
		digitA, digitB = leftPower >= size, rightPower >= size
		if digitA != digitB {
			break
		}
		power++
		if digitA {
			leftPower -= size
			rightPower -= size
		}
		leftPower <<= 1
		rightPower <<= 1
	}
	return power
}

func powersort(p sort.Interface, length int) {
	var top, leftA, rightA, leftB, rightB, power, i int
	if length < 2 {
		return
	}
	rightA = findOrCreateRun(p, 0, length-1)
	if rightA == length-1 {
		return
	}
	// Entries in the runs are indexed by their node power.
	// Node power will never be more than log2(totalLen) + 1
	// Therefore, the runstack needs to contain (log2(totalLen)+1) entries
	runs := make([][2]int, log2(uint64(length))+1)
	entryAt := func(i int) bool { return !(runs[i][0] == 0 && runs[i][1] == 0) }

	for rightA < length-1 {
		leftB, rightB = rightA+1, findOrCreateRun(p, rightA+1, length-1)
		power = nodePower(length, leftA, leftB, rightB)
		for i = top; i > power; i-- {
			if entryAt(i) {
				symMerge(p, runs[i][0], runs[i][1]+1, rightA+1)
				leftA = runs[i][0]
				runs[i][0], runs[i][1] = 0, 0
			}
		}
		runs[power][0], runs[power][1] = leftA, rightA
		top = power
		leftA, rightA = leftB, rightB
	}
	for i = top; i >= 0; i-- {
		if entryAt(i) {
			symMerge(p, runs[i][0], runs[i][1]+1, length)
		}
	}
}

// powersort, an adaptive mergesort, described here: https://arxiv.org/pdf/1805.04154.pdf
// The key improvement powersort makes over the standard library merge sort is in carefully
// choosing when to merge adjacent runs of already sorted data based on their calculated
// position in a nearly-optimal binary merge tree.  It really needs a lower level merge function
// that deals well with unequal merge lengths.  Its behaviour does not appear to play nicely with
// symmerge from the standard libary.
func Powersort(v sort.Interface) {
	powersort(v, v.Len())
}
