// Package sgrad provides facilities for the computation and rendering of
// a finite-difference gradient image from a source Sippimage.
package sgrad

import (
	"image"
    "fmt"
    "math"
)

import (
	. "github.com/Causticity/sipp/simage"
)

// Gradimage stores a gradient image with a complex value at each pixel.
type Gradimage struct {
	// The "pixel" data.
	Pix []complex128
	// The rectangle defining the bounds of the image.
	Rect image.Rectangle
	// The maximum modulus value that occurs in this image. This is useful
	// when computing a histogram of the modulus value.
	MaxMod float64
}

// Use a 2x2 kernel to create a finite-differences gradient image, one pixel
// narrower and shorter than the original. We'd rather reduce the size of the
// output image than arbitrarily wrap around or extend the source image, as
// any such procedure could introduce errors into the statistics.
func Fdgrad(src Sippimage) (grad *Gradimage) {
	// Create the dst image from the bounds of the src
	srect := src.Bounds()
	grad = new(Gradimage)
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
			if modsq > grad.MaxMod {
				grad.MaxMod = modsq
			}
		}
	}
	grad.MaxMod = math.Sqrt(grad.MaxMod)

	return
}

// Render the real and imaginary parts of the gradient image as separate 
// 8-bit grayscale images.
func (grad *Gradimage) Render() (Sippimage, Sippimage) {
	// TODO: Store these in Gradimage and compute them as the image is
	// 		 generated.
	// compute max excursions of the real and imag parts
	var minreal float64 = math.MaxFloat64
	var minimag float64 = math.MaxFloat64
	var maxreal float64 = -math.MaxFloat64
	var maximag float64 = -math.MaxFloat64
	for _, pix := range grad.Pix {
		re := real(pix)
		im := imag(pix)
		if re < minreal {
			minreal = re
		}
		if re > maxreal {
			maxreal = re
		}
		if im < minimag {
			minimag = im
		}
		if im > maximag {
			maximag = im
		}
	}
	// compute scale and offset for each image
	rscale := 255.0 / (maxreal - minreal)
	iscale := 255.0 / (maximag - minimag)
	re := new(SippGray)
	re.Gray = image.NewGray(grad.Rect)
	im := new(SippGray)
	im.Gray = image.NewGray(grad.Rect)
	// scan the complex image, generating the two renderings
	rePix := re.Pix()
	imPix := im.Pix()
	for index, pix := range grad.Pix {
		r := real(pix)
		i := imag(pix)
		rePix[index] = uint8((r - minreal)*rscale)
		imPix[index] = uint8((i - minimag)*iscale)
	}
	return re, im
}