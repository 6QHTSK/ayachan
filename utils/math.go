package utils

// Reciprocal 去除0的倒数为inf的问题
func Reciprocal(num float64) (r float64) {
	if num == 0.0 {
		return r
	}
	return 1.0 / num
}
