package gorl

func IntAbs(i int) uint {
	if i < 0 {
		i *= -1
	}
	return uint(i)
}
