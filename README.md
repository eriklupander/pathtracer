# pathtracer
Pathtracer written in Go
![example](images/super-highres-refl.png)
## Description
Simple unidirectional pathtracer written just for fun. 

Supports:
* Spheres and Planes
* Diffuse or strictly reflective materials
* Movable camera
* Multithreaded rendering and AVX2 accelerated ray transforms

Based on or inspired by:

* My implementation of "The Ray Tracer Challenge" at https://github.com/eriklupander/rt
* Ray in hemisphere code by Hunter Loftis at https://github.com/hunterloftis/pbr/blob/1ce8b1c067eea7cf7298745d6976ba72ff12dd50/pkg/geom/dir.go
* Mask/accumulated color shading by Sam Lapere at https://raytracey.blogspot.com/2016/11/opencl-path-tracing-tutorial-2-path.html

## Gallery
* Main image: 4096 samples, 1280x960, 3 bounces per sample. Took 2h3m to render on MacBook Pro mid 2014.