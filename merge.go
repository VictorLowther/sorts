package sorts

import "sort"

func StdlibStable(vals []int, a Less) {
	sort.Stable(&sortData{data: vals, less: a})

}
