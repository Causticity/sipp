// Copyright Raul Vera 2015-2016

package scomplex

import (
	"image"
	"math"
)

import (
	. "github.com/Causticity/sipp/simage"
)

// A ComplexInt32Image is an image where each pixel is a ComplexInt32.
type ComplexInt32Image struct {
	// The "pixel" data.
	Pix []ComplexInt32
	// The rectangle defining the bounds of the image.
	Rect image.Rectangle
	// The maximum modulus value that occurs in this image. This is useful
	// when computing a histogram of the modulus value. We retain a float64 for
	// precision.
	MaxMod float64
	// Extreme values found in this image
	MinRe, MaxRe, MinIm, MaxIm int32
}

// FromComplexInt32Array wraps an array of complex numbers in a ComplexInt32Image.
func FromComplexInt32Array(cpx []ComplexInt32, width int) (dst *ComplexInt32Image) {
	dst = new(ComplexInt32Image)
	dst.Pix = cpx
	dst.Rect = image.Rect(0, 0, width, len(cpx)/width)
	dst.MinRe = math.MaxInt32
	dst.MinIm = math.MaxInt32
	dst.MaxRe = -math.MaxInt32
	dst.MaxIm = -math.MaxInt32
	dst.MaxMod = 0.0
	for _, c := range cpx {
		reVal := c.Re
		imVal := c.Im
		modsq := float64(reVal*reVal) + float64(imVal*imVal)
		// store the maximum squared value, then take the root afterwards
		if modsq > dst.MaxMod {
			dst.MaxMod = modsq
		}
		if reVal < dst.MinRe {
			dst.MinRe = reVal
		}
		if reVal > dst.MaxRe {
			dst.MaxRe = reVal
		}
		if imVal < dst.MinIm {
			dst.MinIm = imVal
		}
		if imVal > dst.MaxIm {
			dst.MaxIm = imVal
		}
	}
	dst.MaxMod = math.Sqrt(dst.MaxMod)
	return
}

// ToShiftedInt32Complex converts the input image into a ComplexInt32Image,
// multiplying each pixel by (-1)^(x+y), in order for a subsequent FFT to be
// centred properly.
func ToShiftedComplexInt32(src SippImage) (dst *ComplexInt32Image) {
	dst = new(ComplexInt32Image)
	dst.Rect = src.Bounds()
	width := dst.Rect.Dx()
	height := dst.Rect.Dy()
	size := width * height
	dst.Pix = make([]ComplexInt32, size)
	dst.MinRe = math.MaxInt32
	dst.MinIm = 0
	dst.MaxRe = math.MinInt32
	dst.MaxIm = 0
	dst.MaxMod = 0.0
	// Multiply by (-1)^(x+y) while converting the pixels to complex numbers
	var shiftStart int32 = 1
	shift := shiftStart
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := src.IntVal(x, y)*shift
			dst.Pix[i] = ComplexInt32{val, 0}
			shift = -shift
			modsq := float64(val*val)
			if modsq > dst.MaxMod {
				dst.MaxMod = modsq
			}
			if val < dst.MinRe {
				dst.MinRe = val
			}
			if val > dst.MaxRe {
				dst.MaxRe = val
			}
			i++
		}
		shiftStart = -shiftStart
		shift = shiftStart
	}
	dst.MaxMod = math.Sqrt(dst.MaxMod)

	return
}

// Render renders the real and imaginary parts of the image as separate 8-bit
// grayscale images.
func (comp *ComplexInt32Image) Render() (SippImage, SippImage) {
	// compute scale and offset for each image
	rdiv := comp.MaxRe - comp.MinRe
	if rdiv < 1 {
		rdiv = 1
	}
	idiv := comp.MaxIm - comp.MinIm
	if idiv < 1 {
		idiv = 1
	}
	rscale := 255.0 / float64(rdiv)
	iscale := 255.0 / float64(idiv)
	re := new(SippGray)
	re.Gray = image.NewGray(comp.Rect)
	im := new(SippGray)
	im.Gray = image.NewGray(comp.Rect)
	// scan the complex image, generating the two renderings
	rePix := re.Pix()
	imPix := im.Pix()
	for index, pix := range comp.Pix {
		r := pix.Re
		i := pix.Im
		rePix[index] = uint8(float64((r - comp.MinRe)) * rscale)
		imPix[index] = uint8(float64((i - comp.MinIm)) * iscale)
	}
	return re, im
}
