package tracer

import (
	camera2 "github.com/eriklupander/pt/internal/app/camera"
	canvas2 "github.com/eriklupander/pt/internal/app/canvas"
	"github.com/eriklupander/pt/internal/app/geom"
	"github.com/eriklupander/pt/internal/app/scenes"
	"github.com/eriklupander/pt/internal/app/shapes"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const samples = 512
const maxBounces = 8

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

// trace traces the camera ray, collecting colors and emission on each bounce.
// mask/accum based on https://raytracey.blogspot.com/2016/11/opencl-path-tracing-tutorial-2-path.html
func (ctx *Ctx) trace(ray geom.Ray) geom.Tuple4 {
	incrSampleCount()

	ctx.resetMaskAndAccumulatedColors()

	var transformedRay geom.Ray

	for bounces := 0; bounces < maxBounces; bounces++ {
		var hitIndex = 0
		t := 0.0
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
			return geom.Add(ctx.accumColor, geom.Hadamard(ctx.mask, backgroundColor))
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

		// get hit.
		intersection := ctx.intersections[hitIndex]

		// Get point of hit, store in ctx to reuse memory
		PositionPtr(ray, t, &ctx.hitpoint)

		// Use the computations func from Ray-tracer challenge impl that computes surface normal, reflVec etc.
		PrepareComputationForIntersectionPtr(intersection, ray, &ctx.comps, ctx.intersections...)

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

		// DIFFUSE - compute random ray in hemisphere and "overwrite" ray for next iteration
		if intersection.S.GetMaterial().Reflectivity < 1.0 {
			newdir := ctx.RandomConeInHemisphere(ctx.comps.NormalVec, 1.0)
			ray = geom.NewRay(ctx.comps.OverPoint, newdir)

			// First, ADD current color with the hadamard of the current mask and the emission properties of the hit object.
			ctx.accumColor = geom.Add(ctx.accumColor, geom.Hadamard(ctx.mask, intersection.S.GetMaterial().Emission))

			// TODO check if we should terminate if we've hit a light source. No need to keep bouncing

			// The updated mask is used on _the next_ bounce
			// the mask colour picks up surface colours at each bounce
			ctx.mask = geom.Hadamard(ctx.mask, intersection.S.GetMaterial().Color)

			// perform cosine-weighted importance sampling for diffuse surfaces
			ctx.mask = geom.MultiplyByScalar(ctx.mask, geom.Dot(newdir, ctx.comps.NormalVec))
		} else {
			// Full reflection, looks ok-ish but I think the emission is wrong
			ray = geom.NewRay(ctx.comps.OverPoint, ctx.comps.ReflectVec)
			ctx.mask = geom.Hadamard(ctx.mask, intersection.S.GetMaterial().Color)
		}
	}
	return ctx.accumColor
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

// Based on Hunter Loftis path tracer:
// https://github.com/hunterloftis/pbr/blob/1ce8b1c067eea7cf7298745d6976ba72ff12dd50/pkg/geom/dir.go
//
// Cone returns a random vector within a cone about Direction a.
// size is 0-1, where 0 is the original vector and 1 is anything within the original hemisphere.
// https://github.com/fogleman/pt/blob/69e74a07b0af72f1601c64120a866d9a5f432e2f/pt/util.go#L24
func (ctx *Ctx) RandomConeInHemisphere(startVec geom.Tuple4, size float64) geom.Tuple4 {
	u := ctx.rnd.Float64()
	v := ctx.rnd.Float64()
	theta := size * 0.5 * math.Pi * (1 - (2 * math.Acos(u) / math.Pi))
	m1 := math.Sin(theta)
	m2 := math.Cos(theta)
	a2 := v * 2 * math.Pi

	// q should be possible to store in Context?
	q := geom.Tuple4{}
	ctx.RandDirection(&q)

	// should be possible to move s and t into Context?
	s := geom.Tuple4{}
	t := geom.Tuple4{}
	geom.Cross2(&startVec, &q, &s)
	geom.Cross2(&startVec, &s, &t)

	d := geom.Tuple4{}
	geom.MultiplyByScalarPtr(s, m1*math.Cos(a2), &d)
	d = geom.Add(d, geom.MultiplyByScalar(t, m1*math.Sin(a2)))
	d = geom.Add(d, geom.MultiplyByScalar(startVec, m2))
	return geom.Normalize(d)
}

// RandDirection returns a random unit vector (a point on the edge of a unit sphere).
func (ctx *Ctx) RandDirection(out *geom.Tuple4) {
	AngleDirection(ctx.rnd.Float64()*math.Pi*2, math.Asin(ctx.rnd.Float64()*2-1), out)
}
func AngleDirection(theta, phi float64, out *geom.Tuple4) {
	out[0] = math.Cos(theta) * math.Cos(phi)
	out[1] = math.Sin(phi)
	out[2] = math.Sin(theta) * math.Cos(phi)
}
