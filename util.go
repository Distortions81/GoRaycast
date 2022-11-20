package main

import "math"

func distance(a, b xycord) float64 {
	return math.Hypot(float64(a.y-b.x), float64(a.y-b.y))
}

func fixRad(rad float64) float64 {
	for rad < 0 {
		rad += twoPi
	}
	for rad >= twoPi {
		rad -= twoPi
	}
	return rad
}

func degToRad(deg float64) float64 {
	return fixRad(deg * onePi / 180.0)
}
