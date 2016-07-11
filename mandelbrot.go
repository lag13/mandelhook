package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/cmplx"
	"os"
)

const (
	maxEsc          = 3000
	rMin            = -2
	rMax            = .5
	iMin            = -1
	iMax            = 1
	mandelbrotWidth = 75
	gridCellSideLen = 4
)

var (
	scale            = mandelbrotWidth / (rMax - rMin)
	mandelbrotHeight = int(scale * (iMax - iMin))
)

func mandelbrot(a complex128) float64 {
	i := 0
	for z := a; cmplx.Abs(z) < 2 && i < maxEsc; i++ {
		z = z*z + a
	}
	return float64(maxEsc-i) / maxEsc
}

func main() {
	width := mandelbrotWidth * gridCellSideLen
	height := mandelbrotHeight * gridCellSideLen
	bounds := image.Rect(0, 0, width, height)
	b := image.NewNRGBA(bounds)
	draw.Draw(b, bounds, image.NewUniform(color.White), image.ZP, draw.Src)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x%gridCellSideLen == 0 || y%gridCellSideLen == 0 {
				b.Set(x, y, color.Gray{225})
			}
			if x%gridCellSideLen == 0 && x/gridCellSideLen%10 == 0 || y%gridCellSideLen == 0 && y/gridCellSideLen%10 == 0 {
				b.Set(x, y, color.Gray{200})
			}
		}
	}
	mandelbrotSet := buildMandelbrot()
	mx := 0
	for x := gridCellSideLen / 2; x < gridCellSideLen*mandelbrotWidth; x += gridCellSideLen {
		my := 0
		for y := gridCellSideLen / 2; y < gridCellSideLen*mandelbrotHeight; y += gridCellSideLen {
			b.Set(x, y, mandelbrotSet[mx+my*mandelbrotWidth])
			my++
		}
		mx++
	}
	f, err := os.Create("mandelbrot.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = png.Encode(f, b); err != nil {
		fmt.Println(err)
	}
	if err = f.Close(); err != nil {
		fmt.Println(err)
	}
}

func buildMandelbrot() []color.Color {
	mandelbrotSet := make([]color.Color, mandelbrotWidth*mandelbrotHeight)
	for x := 0; x < mandelbrotWidth; x++ {
		for y := 0; y < mandelbrotHeight; y++ {
			fEsc := mandelbrot(complex(float64(x)/scale+rMin, float64(y)/scale+iMin))
			mandelbrotSet[x+y*mandelbrotWidth] = color.Gray{uint8(fEsc * 255)}
		}
	}
	return mandelbrotSet
}
