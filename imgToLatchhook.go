package main

import (
	"flag"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"

	"github.com/lag13/mandelhook/imageutil"
	"github.com/lag13/mandelhook/reducecolors"
)

// TODO: See what the main pain points for this program is and fix them.
// Probably memory allocations?

func main() {
	// TODO: We could just force them to enter a file rather than have this
	// input flag and if they enter nothing then we pull from stdin?
	input := flag.String("input", "mandelbrot.png", "image to load")
	output := flag.String("output", "test.png", "output image name")
	// TODO: Perhaps by default we'll keep all colors which consist of a high
	// enough percentage of the overall picture. Maybe we'll want this behavior
	// and the ability to specify a number.
	numColors := flag.Int("num", 5, "number of unique colors in the final image")
	flag.Parse()
	imgInterface, err := imageutil.LoadImageIntoMemory(*input)
	if err != nil {
		log.Print(err)
		return
	}
	img := reducecolors.PalettedApproach(imgInterface, *numColors)
	// oldImg := newRGBA(img)
	// bounds := img.Bounds()
	// for x := bounds.Min.X; x < bounds.Max.X; x++ {
	// 	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 		neighbors := getValidNeighbors(getAllNeighbors(x, y), bounds)
	// 		// neighbors = getCloseNeighbors(neighbors, img)
	// 		img.Set(x, y, averageNeighboringPixels(neighbors, oldImg))
	// 	}
	// }
	if err := imageutil.WriteImage(*output, img, png.Encode); err != nil {
		log.Print(err)
		return
	}
}

func averageNeighboringPixels(neighbors []image.Point, img image.Image) color.Color {
	var rSum, gSum, bSum uint32
	for _, n := range neighbors {
		r, g, b, _ := img.At(n.X, n.Y).RGBA()
		rSum += r
		gSum += g
		bSum += b
	}
	l := uint32(len(neighbors))
	return color.RGBA{uint8(rSum / l), uint8(gSum / l), uint8(bSum / l), 255}
}

func getAllNeighbors(x int, y int) []image.Point {
	neighbors := []image.Point{}
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			neighbors = append(neighbors, image.Point{x + dx, y + dy})
		}
	}
	return neighbors
}

func getValidNeighbors(neighbors []image.Point, bounds image.Rectangle) []image.Point {
	n := []image.Point{}
	for _, neighbor := range neighbors {
		if bounds.Min.X <= neighbor.X && neighbor.X < bounds.Max.X && bounds.Min.Y <= neighbor.Y && neighbor.Y < bounds.Max.Y {
			n = append(n, neighbor)
		}
	}
	return n
}
