package geom

import "math"

var Epsilon = 0.01

var Up = NewVector(0, 1, 0)
var Right = NewVector(1, 0, 0)

func Eq(a, b float64) bool {
	return math.Abs(a-b) < Epsilon
}
