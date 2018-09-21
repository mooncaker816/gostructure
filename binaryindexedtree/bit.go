package bit

// BIT - Binary Indexed Tree
type BIT []int

// NewBIT returns a binary indexed tree
func NewBIT(a []int) BIT {
	bit := BIT(make([]int, len(a)+1))
	for i, v := range a {
		bit.Add(i, v)
	}
	return bit
}

// Add adds value delta to the elements of BIT which inflect the value of src array
func (bit BIT) Add(srcIdx, delta int) {
	if srcIdx+1 >= len(bit) {
		panic("srcIdx is out of range")
	}
	bitIdx := srcIdx + 1
	for bitIdx < len(bit) {
		bit[bitIdx] += delta
		bitIdx = nextIdx(bitIdx)
	}
}

// Sum returns the prefix sum of src[0:srcIdx+1] (inclusive)
func (bit BIT) Sum(srcIdx int) int {
	if !bit.chkSrcIdx(srcIdx) {
		panic("srcIdx is out of range")
	}
	sum := 0
	bitIdx := srcIdx + 1
	for bitIdx > 0 {
		sum += bit[bitIdx]
		bitIdx = parentIdx(bitIdx)
	}
	return sum
}

// SumRange returns the sum of src[srcLo:srcHi+1] (inclusive)
func (bit BIT) SumRange(srcLo, srcHi int) int {
	if !bit.chkSrcIdx(srcLo) || !bit.chkSrcIdx(srcHi) {
		panic("range is not supported for the underlying array")
	}
	return bit.Sum(srcHi) - bit.Sum(srcLo-1)
}

// last set bit + 1
func nextIdx(i int) int {
	return i + i&(-i)
}

// last set bit - 1 (flip last set bit) <=> i&(i-1)
func parentIdx(i int) int {
	return i - i&(-i)
}

func (bit BIT) chkSrcIdx(i int) bool {
	return i >= 0 && i+1 < len(bit)
}
