package num

import (
	"image"
	"math"
)

// Determines if two numbers are within a certain amount of difference
func within(num1 int, num2 int, rng int) bool {
	return (math.Abs(float64(num2-num1)) <= float64(rng))
}

// Determines if the X and Y coordinates of two points are similar enough to the degree of the range specified
func PntWithin(pnt1 image.Point, pnt2 image.Point, rng int) bool {
	return ((within(pnt1.X, pnt2.X, rng)) && (within(pnt1.Y, pnt2.Y, rng)))
}
