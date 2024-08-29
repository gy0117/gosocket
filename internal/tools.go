package internal

// CeilPow2 2^k >= x，返回最小的2^k
func CeilPow2(n int) int {
	x := 1
	for x < n {
		x *= 2
	}
	return x
}
