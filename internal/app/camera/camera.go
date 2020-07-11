package camera

import (
	"github.com/eriklupander/pt/internal/app/geom"
	"math"
)

type Camera struct {
	Width       int
	Height      int
	Fov         float64
	Transform   geom.Mat4x4
	Inverse     geom.Mat4x4
	PixelSize   float64
	HalfWidth   float64
	HalfHeight  float64
	Aperture    float64
	FocalLength float64
}

func NewCamera(width int, height int, fov float64, from geom.Tuple4, lookAt geom.Tuple4) Camera {
	// Get the length of half the opposite part of the triangle
	halfView := math.Tan(fov / 2)
	aspect := float64(width) / float64(height)
	var halfWidth, halfHeight float64
	if aspect >= 1.0 {
		halfWidth = halfView
		halfHeight = halfView / aspect
	} else {
		halfWidth = halfView * aspect
		halfHeight = halfView
	}
	pixelSize := (halfWidth * 2) / float64(width)

	transform := ViewTransform(from, lookAt, geom.NewVector(0, 1, 0))
	inverse := geom.Inverse(transform)
	return Camera{
		Width:      width,
		Height:     height,
		Fov:        fov,
		Transform:  transform,
		Inverse:    inverse,
		PixelSize:  pixelSize,
		HalfWidth:  halfWidth,
		HalfHeight: halfHeight,
		Aperture:   0.0, // default, pinhole
	}
}

func ViewTransform(from, to, up geom.Tuple4) geom.Mat4x4 {
	// Create a new matrix from the identity matrix.
	vt := geom.IdentityMatrix //Mat4x4{Elems: make([]float64, 16)}
	//copy(vt.Elems, IdentityMatrix.Elems)

	// Sub creates the initial vector between the eye and what we're looking at.
	forward := geom.Normalize(geom.Sub(to, from))

	// Normalize the up vector
	upN := geom.Normalize(up)

	// Use the cross product to get the "third" axis (in this case, not the forward or up one)
	left := geom.Cross(forward, upN)

	// Again, use cross product between the just computed left and forward to get the "true" up.
	trueUp := geom.Cross(left, forward)

	// copy each axis into the matrix
	vt[0] = left.Get(0)
	vt[1] = left.Get(1)
	vt[2] = left.Get(2)

	vt[4] = trueUp.Get(0)
	vt[5] = trueUp.Get(1)
	vt[6] = trueUp.Get(2)

	vt[8] = -forward.Get(0)
	vt[9] = -forward.Get(1)
	vt[10] = -forward.Get(2)

	// finally, move the view matrix opposite the camera position to emulate that the camera has moved.
	return geom.Multiply(vt, geom.Translate(-from.Get(0), -from.Get(1), -from.Get(2)))
}
