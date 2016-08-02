// Package latchhook transforms an image A into another image B representing
// the latch hook diagram for A.
package latchhook

// TODO: Make a latchook out of the gopher and send it to renee freench.

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/lag13/mandelhook/imageutil"
)

// defaultGridCellSideLen is the default number of pixels on the side of a grid cell
// in the latch hook diagram.
const defaultGridCellSideLen = 4

// CreateDiagram creates the latch hook diagram.
func CreateDiagram(img image.Image) image.Image {
	return createDiagram(img, defaultGridCellSideLen)
}

func createDiagram(img image.Image, cellSideLen int) image.Image {
	diagram := skeletonDiagram(img.Bounds(), defaultGridCellSideLen)
	drawDirections(diagram, img, cellSideLen)
	return diagram
}

// skeletonDiagram creates a white image with grid lines on it.
func skeletonDiagram(b image.Rectangle, cellSideLen int) draw.Image {
	skeleton := image.NewRGBA(image.Rect(0, 0, b.Dx()*cellSideLen, b.Dy()*cellSideLen))
	imageutil.FillWithColor(skeleton, color.White)
	drawGridLines(skeleton, cellSideLen)
	return skeleton
}

func drawGridLines(img draw.Image, cellSideLen int) {
	b := img.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			if x%cellSideLen == 0 && x/cellSideLen%10 == 0 || y%cellSideLen == 0 && y/cellSideLen%10 == 0 {
				img.Set(x, y, color.Gray{200})
			} else if x%cellSideLen == 0 || y%cellSideLen == 0 {
				img.Set(x, y, color.Gray{225})
			}
		}
	}
}

// drawDirections puts a symbol into the center of each of the cells. These are
// the directions for how to make the latch hook.
func drawDirections(dst draw.Image, img image.Image, cellSideLen int) {
	b := img.Bounds()
	mx := 0
	for x := cellSideLen / 2; x < cellSideLen*b.Dx(); x += cellSideLen {
		my := 0
		for y := cellSideLen / 2; y < cellSideLen*b.Dy(); y += cellSideLen {
			dst.Set(x, y, img.At(mx, my))
			my++
		}
		mx++
	}
}
