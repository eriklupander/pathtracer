package tracer

import (
	camera2 "github.com/eriklupander/pathtracer/internal/app/camera"
	canvas2 "github.com/eriklupander/pathtracer/internal/app/canvas"
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/scenes"
	"github.com/eriklupander/pathtracer/internal/app/shapes"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const samples = 2048
const maxBounces = 4

var backgroundColor = geom.NewColor(0.15, 0.15, 0.25)

type Ctx struct {
	Id       int
	scene    *scenes.Scene
	canvas   *canvas2.Canvas
	camera   camera2.Camera
	jobsChan chan *job
	wg       *sync.WaitGroup

	rnd *rand.Rand

	// local storage
	pointInView   geom.Tuple4
	pixel         geom.Tuple4
	origin        geom.Tuple4
	direction     geom.Tuple4
	subVec        geom.Tuple4
	comps         Computation
	intersections shapes.Intersections
	hitpoint      geom.Tuple4

	mask       geom.Tuple4
	accumColor geom.Tuple4
}

func NewCtx(id int, scene *scenes.Scene, canvas *canvas2.Canvas, jobsChan chan *job, wg *sync.WaitGroup) *Ctx {
	return &Ctx{Id: id, scene: scene, canvas: canvas, camera: scene.Camera, jobsChan: jobsChan, wg: wg, rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
		// allocate memory
		pointInView: geom.NewPoint(0, 0, -1.0),
		pixel:       geom.NewColor(0, 0, 0),
		origin:      geom.NewPoint(0, 0, 0),
		direction:   geom.NewVector(0, 0, 0),
		subVec:      geom.NewVector(0, 0, 0),
	}
}

func (ctx *Ctx) workerFuncPerLine() {
	for job := range ctx.jobsChan {
		for i := 0; i < ctx.canvas.W; i++ {
			job.col = i
			ctx.renderPixelPathTracer(job.col, job.row)
		}
		ctx.wg.Done()
	}
}

func (ctx *Ctx) renderPixelPathTracer(x, y int) {
	// shade a single pixel
	var cameraRay = geom.NewEmptyRay()
	var finalColor = geom.NewTuple()

	for i := 0; i < samples; i++ {
		// create a ray through the image plane pixel we're rendering. Each ray is randomly offset somewhere in the pixel.
		ctx.rayForPixelPathTracer(x, y, &cameraRay)

		// call the trace function to compute the pixel color for this particular sample
		clr := ctx.trace(cameraRay)

		// add the color to the "final" color
		finalColor = geom.Add(finalColor, clr)
	}

	// Write final pixel color to canvas after dividing finalColor by number of samples to get the avg.
	ctx.canvas.WritePixelMutex(x, y, geom.DivideByScalar(finalColor, float64(samples)))
}

type bounce struct {
	point           geom.Tuple4
	cos             float64
	color           geom.Tuple4
	emission        geom.Tuple4
	diffuse         bool
	refractiveIndex float64
}

// trace traces the camera ray, collecting colors and emission on each bounce.
// mask/accum based on https://raytracey.blogspot.com/2016/11/opencl-path-tracing-tutorial-2-path.html
func (ctx *Ctx) trace(ray geom.Ray) geom.Tuple4 {
	incrSampleCount()
	ctx.resetMaskAndAccumulatedColors()

	var transformedRay geom.Ray
	var bounces = make([]bounce, 0)

	for i := 0; i < maxBounces; i++ {
		ctx.comps = NewComputation()
		intersection, found := ctx.findIntersection(ray, transformedRay)
		if !found {
			return geom.Add(ctx.accumColor, geom.Hadamard(ctx.mask, backgroundColor))
		}

		b := bounce{
			point:    ctx.comps.OverPoint,
			color:    intersection.S.GetMaterial().Color,
			emission: intersection.S.GetMaterial().Emission,
		}

		if intersection.S.GetMaterial().RefractiveIndex != 1.0 {

			// some rays could reflect rather than refract
			if intersection.S.GetMaterial().Reflectivity > ctx.rnd.Float64() {
				ray = geom.NewRay(ctx.comps.OverPoint, ctx.comps.ReflectVec)
				b.diffuse = false
			} else {
				// Find the ratio of first index of refraction to the second.
				nRatio := ctx.comps.N1 / ctx.comps.N2

				// cos(theta_i) is the same as the dot product of the two vectors
				cosI := geom.Dot(ctx.comps.EyeVec, ctx.comps.NormalVec)

				// Find sin(theta_t)^2 via trigonometric identity
				sin2Theta := (nRatio * nRatio) * (1.0 - (cosI * cosI))
				if sin2Theta > 1.0 {
					return black
				}

				// Find cos(theta_t) via trigonometric identity
				cosTheta := math.Sqrt(1.0 - sin2Theta)

				// Compute the direction of the refracted ray
				newdir := geom.Sub(
					geom.MultiplyByScalar(ctx.comps.NormalVec, (nRatio*cosI)-cosTheta),
					geom.MultiplyByScalar(ctx.comps.EyeVec, nRatio))

				ray = geom.NewRay(ctx.comps.UnderPoint, newdir)
				b.refractiveIndex = intersection.S.GetMaterial().RefractiveIndex
				b.diffuse = false
			}

		} else if intersection.S.GetMaterial().Reflectivity < 1.0 {
			// DIFFUSE - compute random ray in hemisphere and "overwrite" ray for next iteration
			b.diffuse = true
			var newdir geom.Tuple4
			ctx.randomVectorInHemisphere(ctx.comps.NormalVec, ctx.rnd, &newdir)
			ray = geom.NewRay(ctx.comps.OverPoint, newdir)
			b.cos = geom.Dot(newdir, ctx.comps.NormalVec)
		} else {
			b.diffuse = false
			// Full reflection, looks ok-ish but I think the emission is wrong
			ray = geom.NewRay(ctx.comps.OverPoint, ctx.comps.ReflectVec)
		}
		bounces = append(bounces, b)
	}
	return ctx.computeColor(bounces)
}

// Time to compute color and light propagation
// Let's try to explain this shit to myself...
// In order for a ray to contribute color to the pixel (for this sample),
// the ray must end up bouncing into a light source. Otherwise, the
// accumulated color will always remain black since mask * emission will be 0,0,0 with non-emissive materials.
// However, let's say a ray actually bounces into a light on it's third bounce.
// Camera -> Mat 1 (cos 0.5) -> Mat 2 (0.75) -> Light.
// Mat 1 is pure reddish (1, 0.5, 0.5) and Mat 2 is blue-ish (0.5, 0.5, 1.0). The light has 5,5,5 as emission.
// This actually means that on first hit, we'll record mask 1,1,1 and accColor

// Hit 1: color: 0,0,0, mask: 1,1,1 * 1,0.5,0.5 ==> 1, 0.5, 0.5.
//        Final mask with cos 0.5 ==> 0.5, 0.25, 0.25

// Hit 2: color 0,0,0, mask: 0.5, 0.25, 0.25 * 0.5, 0.5, 1.0 == 0.25, 0.125, 0.5.
//        Final mask with cos 0.75 ==> 0.1875, 0.06125, 0.375

// Hit light: color 5,5,5 * 0.1875, 0.06125, 0.375 ==> 0.9375, 0.30625, 1,875
//            since it's the last hit (and lights have to color, only emission), we can ignore
//            the last mask.

// So, what we're basically are doing is that we're collecting colors (mask) on each bounce
// multiplied with the cos between the outgoing new vector and the surface's normal vector.
func (ctx *Ctx) computeColor(bounces []bounce) geom.Tuple4 {
	for _, b := range bounces {
		if b.diffuse {
			// First, ADD current color with the hadamard of the current mask and the emission properties of the hit object.
			ctx.accumColor = geom.Add(ctx.accumColor, geom.Hadamard(ctx.mask, b.emission))

			// The updated mask is used on _the next_ bounce
			// the mask colour picks up surface colours at each bounce
			ctx.mask = geom.Hadamard(ctx.mask, b.color)

			// perform cosine-weighted importance sampling for diffuse surfaces
			ctx.mask = geom.MultiplyByScalar(ctx.mask, b.cos)
		} else if b.refractiveIndex != 1.0 {
			// If we have a transparent material, we should kind of "ignore" the hit on the transparent material
			// and instead use whatever color the next bounce has

		} else {
			// the line below propagates reflected light onto the next surface, but turns the sphere
			// into a...sun?
			//ctx.accumColor = geom.Add(ctx.accumColor, geom.Hadamard(ctx.mask, geom.NewColor(9,8,6)))
			ctx.mask = geom.Hadamard(ctx.mask, b.color)
		}
	}
	return ctx.accumColor
}

func (ctx *Ctx) findIntersection(ray geom.Ray, transformedRay geom.Ray) (shapes.Intersection, bool) {
	var hitIndex = -1
	//t := 0.0
	ctx.intersections = ctx.intersections[:0]

	// find all intersections
	for i := range ctx.scene.Objects {

		// transforming the ray into object space makes the intersection math easier
		TransformRayPtr(ray, ctx.scene.Objects[i].GetInverse(), &transformedRay)

		// Call the intersect function provided by the shape implementation (e.g. Sphere, Plane osv)
		// which appends any intersections to the intersection list
		ctx.scene.Objects[i].IntersectLocal(transformedRay, &ctx.intersections)
	}

	// If there were no intersection
	if len(ctx.intersections) == 0 {
		// from https://github.com/straaljager/OpenCL-path-tracing-tutorial-2-Part-2-Path-tracing-spheres/blob/master/opencl_kernel.cl
		return shapes.Intersection{}, false
	}

	// sort intersections if necessary
	if len(ctx.intersections) > 1 {
		sort.Sort(ctx.intersections)
	}

	// loop over all intersections and find the first positive one
	for i := 0; i < len(ctx.intersections); i++ {

		// Check is positive (in front of camera)
		if ctx.intersections[i].T > 0.0 {
			hitIndex = i
			break
		}
	}
	if hitIndex == -1 {
		return shapes.Intersection{}, false
	}
	// get hit.
	intersection := ctx.intersections[hitIndex]

	// Get point of hit, store in ctx to reuse memory
	PositionPtr(ray, intersection.T, &ctx.hitpoint)

	// Use the computations func from Ray-tracer challenge impl that computes surface normal, reflVec etc.
	PrepareComputationForIntersectionPtr(intersection, ray, &ctx.comps, ctx.intersections...)
	return intersection, true
}

func (ctx *Ctx) resetMaskAndAccumulatedColors() {
	ctx.accumColor[0] = 0
	ctx.accumColor[1] = 0
	ctx.accumColor[2] = 0
	ctx.mask[0] = 1
	ctx.mask[1] = 1
	ctx.mask[2] = 1
}

func (ctx *Ctx) rayForPixelPathTracer(x, y int, out *geom.Ray) {

	xOffset := ctx.camera.PixelSize * (float64(x) + ctx.rnd.Float64()) // 0.5
	yOffset := ctx.camera.PixelSize * (float64(y) + ctx.rnd.Float64()) // 0.5

	// this feels a little hacky but actually works.
	worldX := ctx.camera.HalfWidth - xOffset
	worldY := ctx.camera.HalfHeight - yOffset

	ctx.pointInView[0] = worldX
	ctx.pointInView[1] = worldY

	geom.MultiplyByTuplePtr(&ctx.camera.Inverse, &ctx.pointInView, &ctx.pixel)
	geom.MultiplyByTuplePtr(&ctx.camera.Inverse, &originPoint, &ctx.origin)
	geom.SubPtr(ctx.pixel, ctx.origin, &ctx.subVec)
	geom.NormalizePtr(&ctx.subVec, &ctx.direction)

	out.Direction = ctx.direction
	out.Origin = ctx.origin
}

// randomVectorInHemisphere translated into Go from https://raytracey.blogspot.com/2016/11/opencl-path-tracing-tutorial-2-path.html
// The thing is that using this func for diffuse surfaces produces a good and balanced result in the final image,
// while using the randomConeInHemisphere func translated from Hunter Loftis PathTracer produces overexposed highlights.
//
// Need to figure out why.
func (ctx *Ctx) randomVectorInHemisphere(normalVec geom.Tuple4, rnd *rand.Rand, newdir *geom.Tuple4) {
	var rand1 = 2.0 * math.Pi * rnd.Float64()
	var rand2 = rnd.Float64()
	var rand2s = math.Sqrt(rand2)

	/* create a local orthogonal coordinate frame centered at the hitpoint */
	var axis geom.Tuple4
	if math.Abs(normalVec[0]) > 0.1 {
		axis = geom.Up // NewVector(0, 1, 0)
	} else {
		axis = geom.Right //NewVector(1, 0, 0)
	}
	u := geom.Normalize(geom.Cross(axis, normalVec))
	v := geom.Cross(normalVec, u)

	/* use the coordinte frame and random numbers to compute the next ray direction */
	geom.Add3(geom.MultiplyByScalar(u, math.Cos(rand1)*rand2s), geom.MultiplyByScalar(v, math.Sin(rand1)*rand2s), geom.MultiplyByScalar(normalVec, math.Sqrt(1.0-rand2)), newdir)
}

//
//lightBounces := make([]bounce, maxBounces)
//// try some bi-directional path tracing here
//for _, lightSource := range ctx.scene.Objects {
//	if lightSource.GetMaterial().Emission[0] > 0.0 {
//
//		for i:=0;i < maxBounces;i++ {
//			ctx.intersections = ctx.intersections[:0]
//
//			// 1. cast a random ray from the center of the light source into the scene's hemisphere
//			// (This should be replaced by a emitter func later)
//			lightPos := geom.NewPoint(lightSource.GetTransform().Get(0, 3), lightSource.GetTransform().Get(1, 3), lightSource.GetTransform().Get(2, 3))
//			randomVec := ctx.RandomConeInHemisphere(geom.NewVector(0, -1, 0), 0.2)
//			lightRay := geom.NewRay(lightPos, geom.Normalize(randomVec))
//
//			intersection, found := ctx.findIntersection(lightRay, transformedRay)
//			if !found {
//				break
//			}
//			b := bounce{
//				point:    ctx.comps.OverPoint,
//				color:    intersection.S.GetMaterial().Color,
//				emission: intersection.S.GetMaterial().Emission,
//			}
//
//			// DIFFUSE - compute random ray in hemisphere and "overwrite" ray for next iteration
//			if intersection.S.GetMaterial().Reflectivity < 1.0 {
//				b.diffuse = true
//				newdir := ctx.RandomConeInHemisphere(ctx.comps.NormalVec, 1.0)
//				lightRay = geom.NewRay(ctx.comps.OverPoint, newdir)
//				b.cos = geom.Dot(newdir, ctx.comps.NormalVec)
//			} else {
//				b.diffuse = false
//				lightRay = geom.NewRay(ctx.comps.OverPoint, ctx.comps.ReflectVec)
//				b.cos = geom.Dot(ctx.comps.ReflectVec, ctx.comps.NormalVec)
//			}
//			lightBounces[i] = b
//		}
//		// Break after first light source
//		break
//	}
//}
//
//// Now the tricky part starts. We need to cast a "shadow ray" between the last point on each "train".
//p1 := bounces[len(bounces) - 1]
//p2 := lightBounces[len(lightBounces) - 1]
//vec := geom.Sub(p2.point, p1.point)
//length := geom.MagnitudePtr(&vec)
//shadowRay := geom.NewRay(p1.point, vec)
//intersection, found := ctx.findIntersection(shadowRay, transformedRay)
//if found && intersection.T > 0 && intersection.T < length {
//	// Is shadowed, makes no contribution
//} else {
//	// append light bounces to normal bounces in reverse order, including computing stuff...
//	for i:=0;i <len(lightBounces);i++ {
//		bounces = append(bounces, lightBounces[len(lightBounces) - 1-i])
//	}
//}
//revLightBounces := make([]bounce, len(lightBounces))
//for i, j := 0, len(lightBounces)-1; i < j; i, j = i+1, j-1 {
//	revLightBounces[i], revLightBounces[j] = lightBounces[j], lightBounces[i]
//}
//
//bounces = append(bounces, lightBounces...)
