// Package reducecolors contains algorithms which attempt to reduce the number
// of colors in an image so it will be more suitable to make a latch hook
// diagram out of it (can't have too many colors in a latch hook diagram).
package reducecolors

// TODO: Check out this package (https://github.com/lucasb-eyer/go-colorful) seems to
// have some useful tips about programming with colors.

import (
	"image"
	"image/color"
	"sort"

	"github.com/lag13/mandelhook/imageutil"
	"github.com/lucasb-eyer/go-colorful"
)

// PalettedApproach takes an image, constructs a color palette for an image,
// and returns the paletted version of that image. The color palette is simply
// some number of the most frequently occurring colors.
func PalettedApproach(img image.Image, numColors int) *image.Paletted {
	palette := buildPalette(img, numColors)
	b := img.Bounds()
	p := image.NewPaletted(b, palette)
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			p.Set(x, y, img.At(x, y))
		}
	}
	return p
}

// buildPalette returns a slice of the numColors of the most frequently
// occurring colors in img.
func buildPalette(img image.Image, numColors int) []color.Color {
	orderedColors := orderColorsByFrequency(img)
	palette := make([]color.Color, 0, numColors)
	for i := 0; i < numColors; i++ {
		palette = append(palette, orderedColors[i].c)
	}
	return palette
}

// orderColorsByFrequency returns a slice of colorFrequency's sorted by how
// often the color appears.
func orderColorsByFrequency(img image.Image) []colorFrequency {
	colorCounts := getColorCounts(img)
	colors := []colorFrequency{}
	for k, v := range colorCounts {
		colors = append(colors, colorFrequency{k, v})
	}
	sort.Sort(byFrequency(colors))
	return colors
}

type colorFrequency struct {
	c    color.Color
	freq int
}

type byFrequency []colorFrequency

func (f byFrequency) Len() int           { return len(f) }
func (f byFrequency) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f byFrequency) Less(i, j int) bool { return f[i].freq > f[j].freq }

// getColorCounts returns a map where the keys are the colors of the image and
// the values are how many times that color occurred.
func getColorCounts(img image.Image) map[color.Color]int {
	colorCounts := make(map[color.Color]int)
	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			colorCounts[img.At(x, y)]++
		}
	}
	return colorCounts
}

// NaiveFuzzifyImage was my first attempt to average the colors around a
// particular pixel. The original intention for this sort of algorithm was that
// this would make a pixel's color closer to what it "should be". It didn't
// really work out because if the neighborhood around a particular pixel has a
// lot of sharp color contrast (like going from black to white) then the
// resulting color can be a bit unpredictable.
func NaiveFuzzifyImage(img image.Image) image.Image {
	b := img.Bounds()
	fuzzyImg := image.NewRGBA(b)
	imageutil.Convert(fuzzyImg, img)
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			neighborColors := getMooreNeighborhood(x, y, img)
			fuzzyImg.Set(x, y, averageColors(neighborColors))
		}
	}
	return fuzzyImg
}

func averageColors(colors []color.Color) color.Color {
	var rSum, gSum, bSum uint32
	for _, color := range colors {
		r, g, b, _ := color.RGBA()
		rSum += r
		gSum += g
		bSum += b
	}
	l := uint32(len(colors))
	return color.RGBA{uint8(rSum / l), uint8(gSum / l), uint8(bSum / l), 255}
}

func getMooreNeighborhood(x int, y int, img image.Image) []color.Color {
	cs := []color.Color{}
	b := img.Bounds()
	neighbors := getNeighboringPoints(x, y)
	for _, neighbor := range neighbors {
		if b.Min.X <= neighbor.X && neighbor.X < b.Max.X && b.Min.Y <= neighbor.Y && neighbor.Y < b.Max.Y {
			cs = append(cs, img.At(neighbor.X, neighbor.Y))
		}
	}
	return cs
}

func getNeighboringPoints(x int, y int) []image.Point {
	neighbors := []image.Point{}
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			neighbors = append(neighbors, image.Point{x + dx, y + dy})
		}
	}
	return neighbors
}

// LabDistanceApproach tries the same approach as the PalettedApproach function
// (i.e we find the most frequently occurring colors and change all other
// colors into one that is closest to one of those colors) but using the
// CIE-L*a*b color space.
func LabDistanceApproach(img image.Image, numColors int) image.Image {
	palette := buildPalette(img, numColors)
	colorfulPalette := normalizeColors(palette)
	b := img.Bounds()
	newImg := image.NewRGBA(b)
	imageutil.Convert(newImg, img)
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			closestIdx, bestDist := 0, float64(1<<32-1)
			for i, c := range colorfulPalette {
				colorToChange := normalizeColor(img.At(x, y))
				dist := colorToChange.DistanceLab(c)
				if dist < bestDist {
					closestIdx, bestDist = i, dist
				}
			}
			newImg.Set(x, y, palette[closestIdx])
		}
	}
	return newImg
}

func normalizeColors(cs []color.Color) []colorful.Color {
	result := []colorful.Color{}
	for _, c := range cs {
		result = append(result, normalizeColor(c))
	}
	return result
}

func normalizeColor(c color.Color) colorful.Color {
	r, g, b, _ := c.RGBA()
	return colorful.Color{float64(r) / 65535.0, float64(g) / 65535.0, float64(b) / 65535.0}
}
