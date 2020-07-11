package tracer

import "sync"

var raysCast int64
var sampleCount int64

var cLock = sync.Mutex{}

func incrSampleCount() {
	cLock.Lock()
	sampleCount++
	cLock.Unlock()
}
