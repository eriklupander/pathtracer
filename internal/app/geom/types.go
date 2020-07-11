package geom

import "math"

var Epsilon = 0.01

func Eq(a, b float64) bool {
	return math.Abs(a-b) < Epsilon
}
