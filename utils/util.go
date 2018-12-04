package utils

//MinInt 取int中最小的
func MinInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}
