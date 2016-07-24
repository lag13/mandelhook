package imageutil

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"os"
)

// LoadImageIntoMemory loads the specified image file into memory and returns a
// variable representing that image.
func LoadImageIntoMemory(dst draw.Image, file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("could not load file: %v", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %v", err)
	}
	return img, err
}

type encoder func(w io.Writer, m image.Image) error

// WriteImage writes an image to a file with the specified encoding.
func WriteImage(fileName string, img image.Image, encode encoder) (err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("could not create output file %q: %v", fileName, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err != nil {
			err = fmt.Errorf("could not close file %q: %v", fileName, cerr)
		}
	}()
	if err := encode(file, img); err != nil {
		return fmt.Errorf("could not write the output file %q: %v", fileName, err)
	}
	return err
}

// Convert converts the color models of an image from one type to a different
// one. An example use case is converting an RGBA image into grayscale:
//		b := src.Bounds()
//		dst := image.NewGray(image.Bounds())
//		Convert(dst, img)
func Convert(dst draw.Image, src image.Image) {
	b := src.Bounds()
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
}

// FillWithColor fills in an entire image with one uniform color.
func FillWithColor(dst draw.Image, c color.Color) {
	draw.Draw(dst, dst.Bounds(), image.NewUniform(c), image.ZP, draw.Src)
}
