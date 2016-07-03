// An image where each pixel is a complex number.

package scomplex

import (
	"image"
    "fmt"
    "math"
)

import (
	. "github.com/Causticity/sipp/simage"
)

type Compleximage struct {
	// The "pixel" data.
	Pix []complex128
	// The rectangle defining the bounds of the image.
	Rect image.Rectangle
}

// Render the real and imaginary parts of the image as separate 
// 8-bit grayscale images.
func (comp *Compleximage) Render() (SippImage, SippImage) {
	// compute max excursions of the real and imag parts
	var minreal float64 = math.MaxFloat64
	var minimag float64 = math.MaxFloat64
	var maxreal float64 = -math.MaxFloat64
	var maximag float64 = -math.MaxFloat64
	for _, pix := range comp.Pix {
		reVal := real(pix)
		imVal := imag(pix)
		if reVal < minreal {
			minreal = reVal
		}
		if reVal > maxreal {
			maxreal = reVal
		}
		if imVal < minimag {
			minimag = imVal
		}
		if imVal > maximag {
			maximag = imVal
		}
	}
	fmt.Println("maxreal:",maxreal,", minreal:",minreal)
	// compute scale and offset for each image
	rscale := 255.0 / (maxreal - minreal)
	iscale := 255.0 / (maximag - minimag)
	re := new(SippGray)
	re.Gray = image.NewGray(comp.Rect)
	im := new(SippGray)
	im.Gray = image.NewGray(comp.Rect)
	// scan the complex image, generating the two renderings
	rePix := re.Pix()
	imPix := im.Pix()
	for index, pix := range comp.Pix {
		r := real(pix)
		i := imag(pix)
		rePix[index] = uint8((r - minreal)*rscale)
		imPix[index] = uint8((i - minimag)*iscale)
	}
	return re, im
}

