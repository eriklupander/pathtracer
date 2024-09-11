//go:build !noasm && !appengine
// +build !noasm,!appengine

package geom

import "unsafe"

//go:noescape
func __DotProduct(v1, v2, result unsafe.Pointer)

func DotProduct(v1 *Tuple4, v2 *Tuple4, result *Tuple4) {
	__DotProduct(unsafe.Pointer(v1), unsafe.Pointer(v2), unsafe.Pointer(result))
}
