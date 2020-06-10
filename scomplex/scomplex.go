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

type ComplexImage struct {
	// The "pixel" data.
	Pix []complex128
	// The rectangle defining the bounds of the image.
	Rect image.Rectangle
}

// ToShiftedComplex converts the input image into a complex image, multiplying
// each pixel by (-1)^(x+y), in order for a subsequent FFT to be centred
// properly.
func ToShiftedComplex(src SippImage) (comp *ComplexImage) {
	comp = new(ComplexImage)
	comp.Rect = src.Bounds()
	width := comp.Rect.Dx()
	height := comp.Rect.Dy()
	size := width*height
	comp.Pix = make([]complex128, size)
	// Multiply by (-1)^(x+y) while converting the pixels to complex numbers
	shiftStart := 1.0
	shift := shiftStart
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			comp.Pix[i] = complex(src.Val(x,y)*shift, 0)
			i++
			shift = -shift
		}
		shiftStart = -shiftStart
		shift = shiftStart
	}
	return
}

// Render renders the real and imaginary parts of the image as separate 8-bit
// grayscale images.
func (comp *ComplexImage) Render() (SippImage, SippImage) {
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

