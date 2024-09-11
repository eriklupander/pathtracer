//go:build !noasm && !appengine
// +build !noasm,!appengine

package geom

import (
	"unsafe"
)

//go:noescape
func __DP(v1, v2, result unsafe.Pointer)

func DP(v1 *Tuple4, v2 *Tuple4) float64 {
	result := float64(0)
	__DP(unsafe.Pointer(v1), unsafe.Pointer(v2), unsafe.Pointer(&result))
	return result
}
