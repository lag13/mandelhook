package main

import (
	"flag"
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
	originalImg, err := imageutil.LoadImageIntoMemory(*input)
	if err != nil {
		log.Print(err)
		return
	}
	// img := reducecolors.NaiveFuzzifyImage(originalImg)
	// img := reducecolors.PalettedApproach(originalImg, *numColors)
	img := reducecolors.LabDistanceApproach(originalImg, *numColors)
	if err := imageutil.WriteImage(*output, img, png.Encode); err != nil {
		log.Print(err)
		return
	}
}
