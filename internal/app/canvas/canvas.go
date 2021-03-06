package canvas

import (
	"fmt"
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"math"
	"strconv"
	"strings"
	"sync"
)

type Canvas struct {
	W        int
	H        int
	MaxIndex int
	Pixels   []geom.Tuple4
}

func NewCanvas(w int, h int) *Canvas {
	pixels := make([]geom.Tuple4, w*h)
	for i, _ := range pixels {
		pixels[i] = geom.NewColor(0, 0, 0)
	}
	return &Canvas{W: w, H: h, Pixels: pixels, MaxIndex: w * h}
}

func (c *Canvas) WritePixel(col, row int, color geom.Tuple4) {
	if row < 0 || col < 0 || row >= c.H || col > c.W {
		fmt.Println("pixel was out of bounds")
		return
	}
	if row*col > c.MaxIndex {
		fmt.Println("pixel was out of max bounds index bounds")
		return
	}
	c.Pixels[c.toIdx(col, row)] = color
}

var mutex = sync.Mutex{}

func (c *Canvas) WritePixelMutex(col, row int, color geom.Tuple4) {
	if row < 0 || col < 0 || row >= c.H || col > c.W {
		fmt.Println("pixel was out of bounds")
		return
	}
	if row*col > c.MaxIndex {
		fmt.Println("pixel was out of max bounds index bounds")
		return
	}
	mutex.Lock()
	c.Pixels[c.toIdx(col, row)] = color
	mutex.Unlock()
}

func (c *Canvas) WritePixelToIndex(idx int, color geom.Tuple4) {
	c.Pixels[idx] = color
}

func (c *Canvas) ColorAt(col, row int) geom.Tuple4 {
	return c.Pixels[c.toIdx(col, row)]
}

func clamp(color geom.Tuple4, buf *strings.Builder, written *int) {
	for i := 0; i < 3; i++ {
		c := color.Get(i) * 255.0
		rounded := math.Round(c)
		if rounded > 255.0 {
			rounded = 255.0
		} else if rounded < 0.0 {
			rounded = 0.0
		}
		if *written+3 > 69 {
			buf.WriteString("\n")
			buf.WriteString(strconv.Itoa(int(rounded)))
			buf.WriteString(" ")
			*written = 3
		} else {
			buf.WriteString(strconv.Itoa(int(rounded)))
			if *written+6 < 69 {
				buf.WriteString(" ")
			}
			*written = *written + 4
		}

	}
}

func (c *Canvas) toIdx(x, y int) int {
	return y*c.W + x
}
