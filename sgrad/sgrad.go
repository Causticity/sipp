package sgrad

// Create a gradient image from a source image.

import (
	"image"
	. "github.com/Causticity/sipp/simage"
)

type Gradimage struct {
	Pix []complex64
	Stride int // in pixels, i.e. count of complex64s per scanline
	Rect image.Rectangle
	Mod float64
}

// Use a 2x2 kernel to create a finite-differences gradient image, one pixel
// narrowwer and shorter than the original.
func Fdgrad(src *Sippimage) (grad *Gradimage) {
	// Create the dst image from the bounds of the src
	// set up slices to step through the two images
	return
}