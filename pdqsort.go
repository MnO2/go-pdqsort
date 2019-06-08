package pdqsort

import (
	"math/rand"
	"math/bits"
	"strconv"
)

// A type, typically a collection, that satisfies sort.Interface can be
// sorted by the routines in this package. The methods require that the
// elements of the collection be enumerated by an integer index.
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int

	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool

	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)

	ShiftTail(i, j int)
	ShiftHead(i, j int)
	CyclicSwaps(is, js []int)
}

func partialInsertionSort(data Interface, a, b int) bool {
	const MAX_STEPS = 5
	const SHORTEST_SHIFTING = 50

	len := b - a
	i := 1
	for k := 0; k < MAX_STEPS; k += 1 {
		for i < len && !data.Less(i, i-1) {
			i += 1
		}

		if i == len {
			return true
		}

		if len < SHORTEST_SHIFTING {
			return false
		}

		data.Swap(i-1, i)

		data.ShiftTail(a, i)
		data.ShiftHead(i, b)
	}

	return false
}

func insertionSort(data Interface, a, b int) {
	for i := a + 1; i < b; i++ {
		data.ShiftTail(0, i+1)
	}
}

func siftDown(data Interface, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data.Less(first+child, first+child+1) {
			child++
		}
		if !data.Less(first+root, first+child) {
			return
		}
		data.Swap(first+root, first+child)
		root = child
	}
}

func heapSort(data Interface, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDown(data, lo, i, first)
	}
}

func partitionInBlock(data Interface, a, b, pivot int) int {
	const BLOCK = 128

	l := a
	block_l := BLOCK
	start_l := 0
	end_l := 0
	var offsets_l [BLOCK]int

	r := b
	block_r := BLOCK
	start_r := 0
	end_r := 0
	var offsets_r [BLOCK]int

	for {
		isDone := (r - l) <= 2 * BLOCK

		if isDone {
			rem := r - l
			if start_l < end_l || start_r < end_r {
				rem -= BLOCK
			}

			if start_l < end_l {
				block_r = rem
			} else if start_r < end_r {
				block_l = rem
			} else {
				block_l = rem / 2
				block_r = rem - block_l
			}

			if block_l > BLOCK || block_r > BLOCK {
				panic("block_l > BLOCK or block_r > BLOCK")
			}

			if (r-l) != (block_l + block_r) {
				panic("r-l != block_l + block_r")
			}
		}

		if start_l == end_l {
			start_l = 0
			end_l = 0
			elem := l

			for i := 0; i < block_l; i += 1 {
				offsets_l[end_l] = l + i

				if !data.Less(elem, pivot) {
					end_l += 1
				}

				elem += 1
			}
		}

		if start_r == end_r {
			start_r = 0
			end_r = 0
			elem := r

			for i := 0; i < block_r; i += 1 {
				elem -= 1
				offsets_r[end_r] = r - i - 1

				if data.Less(elem, pivot) {
					end_r += 1
				}
			}
		}

		count := min(end_l-start_l, end_r-start_r)
		if count > 0 {
			data.CyclicSwaps(offsets_l[start_l:(start_l+count)], offsets_r[start_r:(start_r+count)])
			start_l += count
			start_r += count
		}

		if start_l == end_l {
			l += block_l
		}

		if start_r == end_r {
			r -= block_r
		}

		if isDone {
			break
		}
	}

	if start_l < end_l {
		if (r-l) != block_l {
			panic("r-l not equal to block_l")
		}

		for start_l < end_l {
			end_l -= 1
			data.Swap(offsets_l[end_l], r-1)
			r -= 1
		}
		return (r-a)
	} else if start_r < end_r {
		if (r-l) != block_r {
			panic("r-l not equal to block_r")
		}
		
		for start_r < end_r {
			end_r -= 1
			data.Swap(l, offsets_r[end_r])
			l += 1
		}
		return (l-a)
	} else {
		return (l-a)
	}
}

func partition(data Interface, a, b, pivot int) (int, bool) {
	mid, wasPartitioned := func() (int, bool) {
		data.Swap(a, pivot)
		pivot = a

		l := a+1
		r := b
		for l < r && data.Less(l, pivot) {
			l += 1
		}

		for l < r && !data.Less(r-1, pivot) {
			r -= 1
		}

		numOfSmallerThanPivotElems := partitionInBlock(data, l, r, pivot)
		return (l - 1 + numOfSmallerThanPivotElems), (l >= r)
	}()

	data.Swap(a, mid)
	return mid, wasPartitioned
}

func partitionEqual(data Interface, a, b, pivot int) int {
	data.Swap(a, pivot)
	pivot = a

	l := a+1
	r := b
	for {
		for l < r && !data.Less(pivot, l) {
			l += 1
		}

		for l < r && data.Less(pivot, r-1) {
			r -= 1
		}

		if l >= r {
			break
		}

		r -= 1
		data.Swap(l, r)
		l += 1
	}

	return l
}

func breakPatterns(data Interface, a, b int) {
	len := b - a
	if len >= 8 {
		modulus := nextPowerOfTwo(uint(len))
		pos := a + (len / 4 * 2)

		for i := 0; i < 3; i += 1 {
			var gen uint = uint(rand.Int())
			
			other := int(gen & (modulus-1))
			if other >= len {
				other -= len
			}
			other += a

			data.Swap(pos - 1 + i, other)
		}
	}
}

func nextPowerOfTwo(n uint) uint {
	var p uint = 1
	for p < n {
		p *= 2
	}

	return p
}

func reverseRange(data Interface, a, b int) {
	i := a
	j := b - 1
	for i < j {
		data.Swap(i, j)
		i += 1
		j -= 1
	}
}

func choosePivot(data Interface, x, y int) (pivot int, likelySorted bool) {
	const SHORTEST_MEDIAN_OF_MEDIANS = 50
	const MAX_SWAPS = 4 * 3

	len := y - x

	a := len / 4 * 1
	b := len / 4 * 2
	c := len / 4 * 3

	swaps := 0

	if len >= 8 {
		sort2 := func(a, b *int) {
			if data.Less(*b, *a) {
				t := *b
				*b = *a
				*a = t

				swaps += 1
			}
		}

		sort3 := func(a, b, c *int) {
			sort2(a, b)
			sort2(b, c)
			sort2(a, b)
		}

		if len >= SHORTEST_MEDIAN_OF_MEDIANS {
			sortAdajacent := func(a *int) {
				t := *a
				t_minus_one := t-1
				t_plus_one := t+1
				sort3(&t_minus_one, a, &t_plus_one)
			}

			sortAdajacent(&a)
			sortAdajacent(&b)
			sortAdajacent(&c)
		}

		sort3(&a, &b, &c)
	}

	if swaps < MAX_SWAPS {
		return x+b, (swaps == 0)
	} else {
		reverseRange(data, a, b)
		return x+(len - 1 - b), true
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func recurse(data Interface, a, b, pred, limit int) {
	if pred > a {
		panic("pred >= a")
	}

	const MAX_INSERTION = 20

	wasBalanced := true
	// wasPartitioned := true

	for {
		len := b - a
		if len <= MAX_INSERTION {
			insertionSort(data, a, b)
			return
		}

		if limit == 0 {
			heapSort(data, a, b)
			return
		}

		if !wasBalanced {
			breakPatterns(data, a, b)
			limit -= 1
		}

		pivot, _ := choosePivot(data, a, b)
		// if wasBalanced && wasPartitioned && likelySorted {
		// 	if partialInsertionSort(data, a, b) {
		// 		return
		// 	}
		// }

		if pred > 0 {
			if !data.Less(pred, pivot) {
				mid := partitionEqual(data, a, b, pivot)
				a = mid
				continue
			}
		}

		mid, _ := partition(data, a, b, pivot)
		wasBalanced = min(mid-a, len-(mid-a)) >= (len / 8)
		// wasPartitioned = wasP

		left_len := mid - a
		right_len := len - (mid - a) - 1
		if left_len < right_len {
			recurse(data, a, mid, pred, limit)
			a = mid+1
			pred = mid
		} else {
			recurse(data, mid+1, b, mid, limit)
			b = mid
		}
	}
}

func quickSort(data Interface, a, b int) {
	n := data.Len()
	if n == 0 {
		return
	}

	limit := strconv.IntSize * 8 - bits.LeadingZeros(uint(n))
	pred := -1
	recurse(data, a, b, pred, limit)
}

// Sort sorts data.
func Sort(data Interface) {
	n := data.Len()
	quickSort(data, 0, n)
}

// Convenience types for common cases

// IntSlice attaches the methods of Interface to []int, sorting in increasing order.
type IntSlice []int

func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p IntSlice) ShiftTail(a, b int) {
	len := b - a
	if len >= 2 && p[b-1] < p[b-2] {
		tmp := p[b-1]
		p[b-1] = p[b-2]

		i := b-3
		for i >= 0 {
			if !(tmp < p[i]) {
				break
			}

			p[i+1] = p[i]
			i -= 1
		}

		p[i+1] = tmp
	}
}

func (p IntSlice) ShiftHead(a, b int) {
	len := b - a
	if len >= 2 && p[a+1] < p[a] {
		tmp := p[a]
		p[a] = p[a+1]

		i := a+2
		for i < len {
			if !(p[i] < tmp) {
				break
			}

			p[i-1] = p[i]
			i += 1
		}

		p[i-1] = tmp
	}
}

func (p IntSlice) CyclicSwaps(is, js []int) {
	count := len(is)
	tmp := p[is[0]]
	p[is[0]] = p[js[0]]

	for i := 1; i < count; i += 1 {
		p[js[i-1]] = p[is[i]]
		p[is[i]] = p[js[i]]
	}

	p[js[count-1]] = tmp
}

// Sort is a convenience method.
func (p IntSlice) Sort() { Sort(p) }

// Float64Slice attaches the methods of Interface to []float64, sorting in increasing order
// (not-a-number values are treated as less than other values).
type Float64Slice []float64

func (p Float64Slice) Len() int           { return len(p) }
func (p Float64Slice) Less(i, j int) bool { return p[i] < p[j] || isNaN(p[i]) && !isNaN(p[j]) }
func (p Float64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Float64Slice) ShiftTail(a, b int) {
	len := b - a
	if len >= 2 {
		for i := len - 1; i >= 1; i -= 1 {
			if !p.Less(i, i-1) {
				break
			}

			p.Swap(i, i-1)
		}
	}
}

func (p Float64Slice) ShiftHead(a, b int) {
	len := b - a
	if len >= 2 {
		for i := 1; i < len; i += 1 {
			if !p.Less(i, i-1) {
				break
			}

			p.Swap(i, i-1)
		}
	}
}

func (p Float64Slice) CyclicSwaps(is, js []int) {
	count := len(is)
	tmp := p[is[0]]
	p[is[0]] = p[js[0]]

	for i := 1; i < count; i += 1 {
		p[js[i-1]] = p[is[i]]
		p[is[i]] = p[js[i]]
	}

	p[js[count-1]] = tmp
}

// isNaN is a copy of math.IsNaN to avoid a dependency on the math package.
func isNaN(f float64) bool {
	return f != f
}

// Sort is a convenience method.
func (p Float64Slice) Sort() { Sort(p) }

// StringSlice attaches the methods of Interface to []string, sorting in increasing order.
type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p StringSlice) ShiftTail(a, b int) {
	len := b - a
	if len >= 2 {
		for i := len - 1; i >= 1; i -= 1 {
			if !p.Less(i, i-1) {
				break
			}

			p.Swap(i, i-1)
		}
	}
}

func (p StringSlice) ShiftHead(a, b int) {
	len := b - a
	if len >= 2 {
		for i := 1; i < len; i += 1 {
			if !p.Less(i, i-1) {
				break
			}

			p.Swap(i, i-1)
		}
	}
}

func (p StringSlice) CyclicSwaps(is, js []int) {
	count := len(is)
	tmp := p[is[0]]
	p[is[0]] = p[js[0]]

	for i := 1; i < count; i += 1 {
		p[js[i-1]] = p[is[i]]
		p[is[i]] = p[js[i]]
	}

	p[js[count-1]] = tmp
}

// Sort is a convenience method.
func (p StringSlice) Sort() { Sort(p) }

// Convenience wrappers for common cases

// Ints sorts a slice of ints in increasing order.
func Ints(a []int) { Sort(IntSlice(a)) }

// Float64s sorts a slice of float64s in increasing order
// (not-a-number values are treated as less than other values).
func Float64s(a []float64) { Sort(Float64Slice(a)) }

// Strings sorts a slice of strings in increasing order.
func Strings(a []string) { Sort(StringSlice(a)) }

// IntsAreSorted tests whether a slice of ints is sorted in increasing order.
func IntsAreSorted(a []int) bool { return IsSorted(IntSlice(a)) }

// Float64sAreSorted tests whether a slice of float64s is sorted in increasing order
// (not-a-number values are treated as less than other values).
func Float64sAreSorted(a []float64) bool { return IsSorted(Float64Slice(a)) }

// StringsAreSorted tests whether a slice of strings is sorted in increasing order.
func StringsAreSorted(a []string) bool { return IsSorted(StringSlice(a)) }

// IsSorted reports whether data is sorted.
func IsSorted(data Interface) bool {
	n := data.Len()
	for i := n - 1; i > 0; i-- {
		if data.Less(i, i-1) {
			return false
		}
	}
	return true
}
