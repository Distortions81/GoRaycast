package main

import "math"

func distance(a, b xycord) float64 {
	dx := a.x - b.x
	dy := a.y - b.y

	return math.Sqrt(dx*dx + dy*dy)
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
