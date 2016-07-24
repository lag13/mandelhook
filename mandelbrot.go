package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"

	"github.com/lag13/mandelhook/imageutil"
	"github.com/lag13/mandelhook/latchhook"
)

const (
	maxEsc          = 3000
	rMin            = -2
	rMax            = .5
	iMin            = -1
	iMax            = 1
	mandelbrotWidth = 75
)

var (
	scale            = mandelbrotWidth / (rMax - rMin)
	mandelbrotHeight = int(scale * (iMax - iMin))
)

func main() {
	mandelbrotSet := buildMandelbrot()
	b := latchhook.CreateDiagram(mandelbrotSet)
	if err := imageutil.WriteImage("mandelbrot.png", b, png.Encode); err != nil {
		fmt.Println(err)
		return
	}
}

func mandelbrot(a complex128) float64 {
	i := 0
	for z := a; cmplx.Abs(z) < 2 && i < maxEsc; i++ {
		z = z*z + a
	}
	return float64(maxEsc-i) / maxEsc
}

func buildMandelbrot() image.Image {
	mandelbrotSet := image.NewGray(image.Rect(0, 0, mandelbrotWidth, mandelbrotHeight))
	for x := 0; x < mandelbrotWidth; x++ {
		for y := 0; y < mandelbrotHeight; y++ {
			fEsc := mandelbrot(complex(float64(x)/scale+rMin, float64(y)/scale+iMin))
			mandelbrotSet.Set(x, y, color.Gray{uint8(fEsc * 255)})
		}
	}
	return mandelbrotSet
}
