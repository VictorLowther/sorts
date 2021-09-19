package sorts

// Powersort, an adaptive mergesort, described here: https://arxiv.org/pdf/1805.04154.pdf
// The key improvement powersort makes over the standard library merge sort is in carefully
// choosing when to merge adjacent runs of already sorted data based on their calculated
// position in a nearly-optimal binary merge tree.

var (
	minRunLen = 32
)

type powerSort struct {
	sortData
}

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
func (p *powerSort) findOrCreateRun(left, right int) (last int) {
	last = left
	if last == right {
		return
	}
	last++
	if p.less(last, last-1) {
		// Assume we are in a descending run.
		// Stop at the first item that is not in strict descending order.
		// Reverse the list that we found that was in strict descending order.
		for ; last < right; last++ {
			if !p.less(last+1, last) {
				break
			}
		}
		p.reverse(left, last)
	} else {
		// Assume we are in a weakly ascending run.
		// Stop at the first item that is not in weakly ascensing order.
		for ; last < right; last++ {
			if p.less(last+1, last) {
				break
			}
		}
	}

	if last-left < minRunLen-1 {
		// Found a too-short natural run that we can extend.
		// Extend and sort it using an insertion sort.
		right = min(right, left+minRunLen-1)
		p.insertionSort(left, right, last)
		last = right
	}
	return
}

// nodePower is adapted from nodePowerBitwise in
// https://github.com/sebawild/nearly-optimal-mergesort-code/blob/master/src/wildinter/net/mergesort/PowerSort.java
//
// I have not spent enough time figuring out exactly what nodePower does
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
// zeros in an integer.  Go has the former in bits.LeadingZeros, but on
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
func (p *powerSort) nodePower(left, mid, right int) int {
	leftPower := uint(left) + uint(mid)
	rightPower := uint(mid) + uint(right) + 1
	size := uint(len(p.data))
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

func (p *powerSort) sort() {
	var top, leftA, rightA, leftB, rightB, power, i, size int
	size = len(p.data)
	if size < 2 {
		return
	}
	rightA = p.findOrCreateRun(0, size-1)
	if rightA == size-1 {
		return
	}
	// Entries in the runs are indexed by their node power.
	// Node power will never be more than log2(totalLen) + 1
	// Therefore, the runstack needs to contain (log2(totalLen)+1) entries
	runs := make([][2]int, log2(uint64(size))+1)
	entryAt := func(i int) bool { return !(runs[i][0] == 0 && runs[i][1] == 0) }

	for rightA < size-1 {
		leftB, rightB = rightA+1, p.findOrCreateRun(rightA+1, size-1)
		power = p.nodePower(leftA, leftB, rightB)
		for i = top; i > power; i-- {
			if entryAt(i) {
				p.symMerge(runs[i][0], runs[i][1]+1, rightA+1)
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
			p.symMerge(runs[i][0], runs[i][1]+1, size)
		}
	}
}

func Powersort(vals []int, a Less) {
	ps := &powerSort{}
	ps.data = vals
	ps.less = a
	ps.sort()
}
