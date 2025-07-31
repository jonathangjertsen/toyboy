package blip

func FixedOne(bits int) int {
	return 1 << bits
}

func Float64ToFixed(f float64, bits int) int {
	f *= float64(FixedOne(bits))
	f += 0.5
	// original impl calls floor(f), doesn't seem neccessary
	return int(f)
}
