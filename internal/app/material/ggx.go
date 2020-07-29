package material

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"math"
)

// Based on http://www.codinglabs.net/article_physically_based_rendering_cook_torrance.aspx
func chiGGX(v float64) float64 {
	if v > 0 {
		return 1
	}
	return 0
}

func GGXDistribution(n geom.Tuple4, h geom.Tuple4, roughness float64) float64 {
	NoH := geom.Dot(n, h)
	alpha2 := roughness * roughness
	NoH2 := NoH * NoH
	den := NoH2*alpha2 + (1 - NoH2)
	return (chiGGX(NoH) * alpha2) / (math.Pi * den * den)
}
func GGXPartialGeometryTerm(v, n, h geom.Tuple4, alpha float64) float64 {
	VoH2 := saturate(geom.Dot(v, h))
	chi := chiGGX(VoH2 / saturate(geom.Dot(v, n)))
	VoH2 = VoH2 * VoH2
	tan2 := (1 - VoH2) / VoH2
	return (chi * 2) / (1 + math.Sqrt(1+alpha*alpha*tan2))
}
func FresnelSchlick(cosT float64, F0 geom.Tuple4) geom.Tuple4 {
	return geom.MultiplyByScalar(geom.Add(F0, geom.Sub(geom.Tuple4{1, 1, 1, 1}, F0)), math.Pow(1-cosT, 5))
}
func CalcColor(ior float64) geom.Tuple4 {
	// Calculate colour at normal incidence
	F0 := math.Abs((1.0 - ior) / (1.0 + ior))
	F0 = F0 * F0
	//F0 = lerp(F0, materialColour.rgb, metallic);
	return geom.Tuple4{}
}

// from https://docs.microsoft.com/en-us/windows/win32/direct3dhlsl/dx-graphics-hlsl-lerp
func lerp(x, y, s geom.Tuple4) geom.Tuple4 {
	return geom.Add(x, geom.Hadamard(s, geom.Sub(y, x)))
}
func saturate(val float64) float64 {
	if val > 1.0 {
		return 1.0
	}
	if val < 0.0 {
		return 0.0
	}
	return val
}
