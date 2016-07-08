// Copyright Raul Vera 2015-2016

// Package sgrad provides facilities for the computation and rendering of
// a finite-difference gradient image from a source SippImage.
package sgrad

import (
	"image"
    "fmt"
    "math"
)

import (
	. "github.com/Causticity/sipp/simage"
	. "github.com/Causticity/sipp/scomplex"
)

// GradImage stores a gradient image with a complex value at each pixel.
type GradImage struct {
	// A GradImage is a complex image
	ComplexImage
	// The maximum modulus value that occurs in this image. This is useful
	// when computing a histogram of the modulus value.
	MaxMod float64
}

// Use a 2x2 kernel to create a finite-differences gradient image, one pixel
// narrower and shorter than the original. We'd rather reduce the size of the
// output image than arbitrarily wrap around or extend the source image, as
// any such procedure could introduce errors into the statistics.
func Fdgrad(src SippImage) (grad *GradImage) {
	// Create the dst image from the bounds of the src
	srect := src.Bounds()
	grad = new(GradImage)
	grad.Rect = image.Rect(0,0,srect.Dx()-1,srect.Dy()-1)
	grad.Pix = make([]complex128, grad.Rect.Dx()*grad.Rect.Dy())
	grad.MaxMod = 0
	
	fmt.Println("source image rect:<", srect, ">")
	fmt.Println("gradient image rect:<", grad.Rect, ">")
	fmt.Println("Gradient image no. of pixels:<", len(grad.Pix), ">")
	
	// grad[x,y] = complex(src[x+1,y+1] - src[x,y], src[x+1,y]-src[x,y+1])
	dsti := 0
	for y := 0; y < grad.Rect.Dy(); y++ {
		for x := 0; x < grad.Rect.Dx(); x++ {
			re := src.Val(x+1,y+1) - src.Val(x,y)
			im := src.Val(x+1,y) - src.Val(x,y+1)
			grad.Pix[dsti] = complex(re, im)
			dsti++
			modsq := re*re + im*im
			// store the maximum squared value, then take the root afterwards
			if modsq > grad.MaxMod {
				grad.MaxMod = modsq
			}
		}
	}
	grad.MaxMod = math.Sqrt(grad.MaxMod)

	return
}

