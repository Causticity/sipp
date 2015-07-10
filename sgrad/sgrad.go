package sgrad

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

func Fdgrad(src *Sippimage) (grad *Gradimage) {
	// Create the dst image from the bounds of the src
	// set up slices to step through the two images
	return
}