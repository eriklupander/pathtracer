package tracer

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"math"
	"math/rand"
)

// Based on Hunter Loftis path tracer:
// https://github.com/hunterloftis/pbr/blob/1ce8b1c067eea7cf7298745d6976ba72ff12dd50/pkg/geom/dir.go
//
// Cone returns a random vector within a cone about Direction a.
// size is 0-1, where 0 is the original vector and 1 is anything within the original hemisphere.
// https://github.com/fogleman/pt/blob/69e74a07b0af72f1601c64120a866d9a5f432e2f/pt/util.go#L24
func (ctx *Ctx) randomConeInHemisphere(startVec geom.Tuple4, size float64, rnd *rand.Rand) geom.Tuple4 {
	u := rnd.Float64()
	v := rnd.Float64()
	theta := size * 0.5 * math.Pi * (1 - (2 * math.Acos(u) / math.Pi))
	m1 := math.Sin(theta)
	m2 := math.Cos(theta)
	a2 := v * 2 * math.Pi

	// q should be possible to store in Context?
	q := geom.Tuple4{}
	randDirection(&q, rnd)

	// should be possible to move s and t into Context?
	s := geom.Tuple4{}
	t := geom.Tuple4{}
	geom.Cross2(&startVec, &q, &s)
	geom.Cross2(&startVec, &s, &t)

	d := geom.Tuple4{}
	d = geom.MultiplyByScalar(s, m1*math.Cos(a2))
	d = geom.Add(d, geom.MultiplyByScalar(t, m1*math.Sin(a2)))
	d = geom.Add(d, geom.MultiplyByScalar(startVec, m2))
	return geom.Normalize(d)
}

// randDirection returns a random unit vector (a point on the edge of a unit sphere).
func randDirection(out *geom.Tuple4, rnd *rand.Rand) {
	angleDirection(rnd.Float64()*math.Pi*2, math.Asin(rnd.Float64()*2-1), out)
}

func angleDirection(theta, phi float64, out *geom.Tuple4) {
	out[0] = math.Cos(theta) * math.Cos(phi)
	out[1] = math.Sin(phi)
	out[2] = math.Sin(theta) * math.Cos(phi)
}

// Also from Hunter Loftis pathtracer. This one produces a similar, but definitely darker
// result than the one I'm using from
// https://raytracey.blogspot.com/2016/11/opencl-path-tracing-tutorial-2-path.html.
func RandHemi(normalVec geom.Tuple4, rnd *rand.Rand) geom.Tuple4 {
	u := rnd.Float64()
	v := rnd.Float64()
	theta := 2 * math.Pi * u
	phi := math.Acos(2*v - 1)
	x := math.Sin(phi) * math.Cos(theta)
	y := math.Sin(phi) * math.Sin(theta)
	z := math.Cos(phi)
	dir := geom.Normalize(geom.Tuple4{x, y, z, 0})
	if geom.Dot(normalVec, dir) < 0 {
		return geom.Negate(dir)
	}
	return dir
}
