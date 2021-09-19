// sorts is a playground for experimenting with various stable sort algorithms.
// All these conform to the usual Go sort interface.  For now, there is just a basic
// insertion sort, the stable sort from the standard library, and an implementation
// of powersort.
package sorts

import "sort"

// Sort is any function that takes an array of integers and
// a Less function and sorts that array based on it.
type Sort func(sort.Interface)
