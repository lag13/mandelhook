// Package reducecolors contains algorithms which attempt to reduce the number
// of colors in an image so it will be more suitable to make a latch hook
// diagram out of it (can't have too many colors in a latch hook diagram).
package reducecolors

// Check out this package (https://github.com/lucasb-eyer/go-colorful) seems to
// have some useful tips about programming with colors.

import (
	"image"
	"image/color"
	"sort"
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

// buildPalette returns a color.Palette for a particular image.
func buildPalette(img image.Image, numColors int) color.Palette {
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

func getColorsToKeepAndRemove(colorFrequencies []colorFrequency, keepNum int) ([]color.Color, []color.Color) {
	keep := []color.Color{}
	for i := 0; i < keepNum && i < len(colorFrequencies); i++ {
		keep = append(keep, colorFrequencies[i].c)
	}
	remove := []color.Color{}
	for i := keepNum; i < len(colorFrequencies); i++ {
		remove = append(remove, colorFrequencies[i].c)
	}
	return keep, remove
}
