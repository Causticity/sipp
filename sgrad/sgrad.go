// Copyright Raul Vera 2015-2020

// Package sgrad provides facilities for the computation of a finite-difference
// gradient image from a source SippImage and a 2x2 kernel.
// There are two versions, one using float64s and complex128s, and another using
// int32s. The latter makes it easier to guarantee bit accuracy and numerical
// stability. It is not intended as a performance optimisation.
package sgrad

import (
	"image"
	"math"
)

import (
	. "github.com/Causticity/sipp/scomplex"
	. "github.com/Causticity/sipp/simage"
)

// SippGradKernels are 2x2 arrays of complex numbers, defined in the same way as
// images are stored in memory, i.e. in row-major order from the top-left corner
// down.
type SippGradKernel [2][2]complex128
type SippGradInt32Kernel [2][2]ComplexInt32

var defaultKernel = SippGradKernel {
	{-1 + 0i, 0 + 1i},
	{0 - 1i, 1 + 0i},
}

var defaultInt32Kernel = SippGradInt32Kernel {
	{{-1, 0}, {0, 1}},
	{{0, -1}, {1, 0}},
}

// TODO: The non-int32 functions below could be reimplemented to use only
// floating-point arithmetic, with a conversion to complex only at the end. As
// it is now it goes back and forth unnecessarily. This is optimisation and
// should be done only with proper profiling and a specific performance target.

// byKernel applies a SippGradKernel to a pixel and its neighbours to produce a
// finite difference as a complex number.
func byKernel(kern SippGradKernel, pix, right, below, belowRight float64) complex128 {
	return kern[0][0]*complex(pix, 0) +
		kern[0][1]*complex(right, 0) +
		kern[1][0]*complex(below, 0) +
		kern[1][1]*complex(belowRight, 0)
}

// byInt32Kernel applies a SippGradInt32Kernel to a pixel and its neighbours to
// produce a finite difference as a ComplexInt32.
func byInt32Kernel(kern SippGradInt32Kernel,
	               pix, right, below, belowRight int32) ComplexInt32 {
	return kern[0][0].Mult(ComplexInt32{pix, 0}).Add(
		   kern[0][1].Mult(ComplexInt32{right, 0}).Add(
		   kern[1][0].Mult(ComplexInt32{below, 0}).Add(
		   kern[1][1].Mult(ComplexInt32{belowRight, 0}))))
}

// Use a SippGradKernel to create a finite-differences complex gradient image,
// one pixel narrower and shorter than the original. We'd rather reduce the size
// of the output image than arbitrarily wrap around or extend the source image,
// as any such procedure could introduce errors into the statistics.
func FdgradKernel(src SippImage, kern SippGradKernel) (grad *ComplexImage) {
	// Create the dst image from the bounds of the src
	srect := src.Bounds()
	grad = new(ComplexImage)
	grad.Rect = image.Rect(0, 0, srect.Dx()-1, srect.Dy()-1)
	grad.Pix = make([]complex128, grad.Rect.Dx()*grad.Rect.Dy())

	dsti := 0
	for y := 0; y < grad.Rect.Dy(); y++ {
		for x := 0; x < grad.Rect.Dx(); x++ {
			val := byKernel(kern, src.Val(x, y),
				src.Val(x+1, y), src.Val(x, y+1), src.Val(x+1, y+1))
			grad.Pix[dsti] = val
			dsti++
			modsq := real(val)*real(val) + imag(val)*imag(val)
			// store the maximum squared value, then take the root afterwards
			if modsq > grad.MaxMod {
				grad.MaxMod = modsq
			}
		}
	}
	grad.MaxMod = math.Sqrt(grad.MaxMod)

	return
}

// Use a SippGradInt32Kernel to create a finite-differences ComplexInt32
// gradient image, one pixel narrower and shorter than the original. We'd rather
// reduce the size of the output image than arbitrarily wrap around or extend
// the source image, as any such procedure could introduce errors into the
// statistics.
func FdgradInt32Kernel(src SippImage,
					   kern SippGradInt32Kernel) (grad *ComplexInt32Image) {
	// Create the dst image from the bounds of the src
	srect := src.Bounds()
	grad = new(ComplexInt32Image)
	grad.Rect = image.Rect(0, 0, srect.Dx()-1, srect.Dy()-1)
	grad.Pix = make([]ComplexInt32, grad.Rect.Dx()*grad.Rect.Dy())
	grad.MinRe = math.MaxInt32
	grad.MaxRe = math.MinInt32
	grad.MinIm = math.MaxInt32
	grad.MaxIm = math.MinInt32

	dsti := 0
	for y := 0; y < grad.Rect.Dy(); y++ {
		for x := 0; x < grad.Rect.Dx(); x++ {
			val := byInt32Kernel(kern, src.IntVal(x, y),
				src.IntVal(x+1, y), src.IntVal(x, y+1), src.IntVal(x+1, y+1))
			grad.Pix[dsti] = val
			dsti++
			if val.Re < grad.MinRe {
				grad.MinRe = val.Re
			}
			if val.Re > grad.MaxRe {
				grad.MaxRe = val.Re
			}
			if val.Im < grad.MinIm {
				grad.MinIm = val.Im
			}
			if val.Im > grad.MaxIm {
				grad.MaxIm = val.Im
			}
		}
	}

	return
}

// Use a default SippGradKernel to compute a finite-differences gradient. See
// FdgradKernel for details.
func Fdgrad(src SippImage) (grad *ComplexImage) {
	return FdgradKernel(src, defaultKernel)
}

// Use a default SippGradInt32Kernel to compute a finite-differences gradient.
// See FdgradInt32Kernel for details.
func FdgradInt32(src SippImage) (grad *ComplexInt32Image) {
	return FdgradInt32Kernel(src, defaultInt32Kernel)
}
