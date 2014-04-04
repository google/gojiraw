package graphics

func MaxF(x1, x2 float32) float32 {
	if x1 > x2 {
		return x1
	}
	return x2
}

func MinF(x1, x2 float32) float32 {
	if x1 < x2 {
		return x1
	}
	return x2
}

// Returns sin(2 * atan(d))
func Sin2Atan(d float32) float32 {
	return 2.0 * d / (1.0 + d*d)
}

// Returns cos(2 * atan(d))
func Cos2Atan(d float32) float32 {
	return (1.0 - d*d) / (1.0 + d*d)
}
