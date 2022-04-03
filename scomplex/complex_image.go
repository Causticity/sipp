// Copyright Raul Vera

package scomplex

import (
	"image"
	"math"
)

import (
	. "github.com/Causticity/sipp/simage"
)

// A ComplexImage is an image where each pixel is a Go complex128.
type ComplexImage struct {
	// The "pixel" data.
	Pix []complex128
	// The rectangle defining the bounds of the image.
	Rect image.Rectangle
	// The maximum modulus value that occurs in this image. This is useful
	// when computing a histogram of the modulus value.
	MaxMod float64
	// Extreme values found in this image
	MinRe, MaxRe, MinIm, MaxIm float64
}

// FromComplexArray wraps an array of complex numbers in a ComplexImage.
func FromComplexArray(cpx []complex128, width int) (dst *ComplexImage) {
	dst = new(ComplexImage)
	dst.Pix = cpx
	dst.Rect = image.Rect(0, 0, width, len(cpx)/width)
	dst.MinRe = math.MaxFloat64
	dst.MinIm = math.MaxFloat64
	dst.MaxRe = -math.MaxFloat64
	dst.MaxIm = -math.MaxFloat64
	dst.MaxMod = 0.0
	for _, c := range cpx {
		re := real(c)
		im := imag(c)
		modsq := re*re + im*im
		// store the maximum squared value, then take the root afterwards
		if modsq > dst.MaxMod {
			dst.MaxMod = modsq
		}
		if re < dst.MinRe {
			dst.MinRe = re
		}
		if re > dst.MaxRe {
			dst.MaxRe = re
		}
		if im < dst.MinIm {
			dst.MinIm = im
		}
		if im > dst.MaxIm {
			dst.MaxIm = im
		}
	}
	dst.MaxMod = math.Sqrt(dst.MaxMod)

	return
}

// ToShiftedComplex converts the input image into a ComplexImage, multiplying
// each pixel by (-1)^(x+y), in order for a subsequent FFT to be centred
// properly.
func ToShiftedComplex(src SippImage) (dst *ComplexImage) {
	dst = new(ComplexImage)
	dst.Rect = src.Bounds()
	dst.MinRe = math.MaxFloat64
	dst.MinIm = 0
	dst.MaxRe = -math.MaxFloat64
	dst.MaxIm = 0
	dst.MaxMod = 0.0
	width := dst.Rect.Dx()
	height := dst.Rect.Dy()
	size := width * height
	dst.Pix = make([]complex128, size)
	// Multiply by (-1)^(x+y) while converting the pixels to complex numbers
	shiftStart := 1.0
	shift := shiftStart
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := src.Val(x, y) * shift
			dst.Pix[i] = complex(val, 0)
			i++
			shift = -shift
			modsq := val*val
			if modsq > dst.MaxMod {
				dst.MaxMod = modsq
			}
			if val < dst.MinRe {
				dst.MinRe = val
			}
			if val > dst.MaxRe {
				dst.MaxRe = val
			}
		}
		shiftStart = -shiftStart
		shift = shiftStart
	}
	dst.MaxMod = math.Sqrt(dst.MaxMod)

	return
}

// Render renders the real and imaginary parts of the image as separate 8-bit
// grayscale images.
func (comp *ComplexImage) Render() (SippImage, SippImage) {
	// compute scale and offset for each image
	rdiv := comp.MaxRe - comp.MinRe
	if rdiv < 1.0 {
		rdiv = 1.0
	}
	idiv := comp.MaxIm - comp.MinIm
	if idiv < 1.0 {
		idiv = 1.0
	}
	rscale := 255.0 / rdiv
	iscale := 255.0 / idiv
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
		rePix[index] = uint8((r - comp.MinRe) * rscale)
		imPix[index] = uint8((i - comp.MinIm) * iscale)
	}
	return re, im
}
