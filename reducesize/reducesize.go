// Package reducesize provides algorithms which reduces the size of an image so
// it is more suitable to make a latch hook diagram out of (a latch hook with
// lots of cells takes a loooong time).
package reducesize

import (
	"image"
	"image/color"
)

// NaiveReduceImageSize (so named because it was the first one I came up with)
// tries to map every pixel from a bigger image into a pixel in the smaller
// image of given width and height. If multiple pixels from the bigger image
// get mapped to the same pixel in the smaller image then those pixels are
// averaged together.
func NaiveReduceImageSize(newWidth int, newHeight int, img image.Image) image.Image {
	b := img.Bounds()
	xRatio := float64(newWidth) / float64(b.Dx())
	yRatio := float64(newHeight) / float64(b.Dy())
	m := make(map[image.Point][]color.Color)
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			nx := int(float64(x) * xRatio)
			ny := int(float64(y) * yRatio)
			p := image.Point{nx, ny}
			m[p] = append(m[p], img.At(x, y))
		}
	}
	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	for p, cs := range m {
		resized.Set(p.X, p.Y, avgColors(cs))
	}
	return resized
}

// avgColors takes a slice of colors and averages them together.
func avgColors(cs []color.Color) color.Color {
	var r, g, b, a uint32
	for _, c := range cs {
		rTemp, gTemp, bTemp, aTemp := c.RGBA()
		r += rTemp
		g += gTemp
		b += bTemp
		a += aTemp
	}
	l := uint32(len(cs))
	return color.RGBA{uint8(r / l), uint8(g / l), uint8(b / l), uint8(a / l)}
}
