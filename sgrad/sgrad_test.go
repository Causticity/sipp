// Copyright Raul Vera 2020

// Tests for package sgrad.

package sgrad

import (
	//"fmt"
	"math"
	"reflect"
	"testing"
)

import (
	. "github.com/Causticity/sipp/scomplex"
	. "github.com/Causticity/sipp/simage"
	. "github.com/Causticity/sipp/sipptesting"
	. "github.com/Causticity/sipp/sipptesting/sipptestcore"
)

var smallPicGrad = []complex128{
	5 - 3i, 5 - 3i, 5 - 3i,
	5 - 3i, 5 - 3i, 5 - 3i,
	5 - 3i, 5 - 3i, 5 - 3i,
}

var smallPicGradMaxMod = math.Sqrt(34.0)

var identityKernel = SippGradKernel{
	{1 + 0i, 0 + 0i},
	{0 + 0i, 0 + 0i},
}

var imagIdentityKernel = SippGradKernel{
	{0 + 1i, 0 + 0i},
	{0 + 0i, 0 + 0i},
}
var kieransKernel = SippGradKernel{
	{-1 + 0i, 0 - 1i},
	{0 + 1i, 1 + 0i},
}

var smallPicGradKieransKernel = []complex128{
	5 + 3i, 5 + 3i, 5 + 3i,
	5 + 3i, 5 + 3i, 5 + 3i,
	5 + 3i, 5 + 3i, 5 + 3i,
}

var smallPicGradInt32 = []ComplexInt32{
	{5, -3}, {5, -3}, {5, -3},
	{5, -3}, {5, -3}, {5, -3},
	{5, -3}, {5, -3}, {5, -3},
}

var identityKernelInt32 = SippGradInt32Kernel{
	{{1, 0}, {0, 0}},
	{{0, 0}, {0, 0}},
}

var imagIdentityKernelInt32 = SippGradInt32Kernel{
	{{0, 1}, {0, 0}},
	{{0, 0}, {0, 0}},
}
var kieransKernelInt32 = SippGradInt32Kernel{
	{{-1, 0}, {0, -1}},
	{{0, 1}, {1, 0}},
}

var smallPicGradKieransKernelInt32 = []ComplexInt32{
	{5, 3}, {5, 3}, {5, 3},
	{5, 3}, {5, 3}, {5, 3},
	{5, 3}, {5, 3}, {5, 3},
}

var CosxCosyTinyGradInt32 = []ComplexInt32 {
{-23,0}, {-32,-12}, {-37,-22}, {-40,-31}, {-38,-37}, {-31,-39}, {-23,-38}, {-12,-32}, {0,-24}, {13,-13}, {24,0}, {32,12}, {37,23}, {39,32}, {37,37}, {31,39}, {23,38}, {12,32}, {0,23},
{-32,12}, {-37,0}, {-39,-12}, {-37,-23}, {-31,-32}, {-23,-37}, {-13,-39}, {0,-37}, {12,-32}, {23,-23}, {32,-12}, {37,0}, {39,12}, {37,23}, {31,32}, {23,37}, {12,39}, {0,38}, {-12,32},
{-38,23}, {-39,13}, {-38,0}, {-32,-11}, {-24,-22}, {-13,-31}, {0,-37}, {13,-39}, {23,-38}, {32,-32}, {38,-23}, {39,-13}, {37,0}, {31,12}, {23,23}, {12,32}, {0,37}, {-12,39}, {-23,38},
{-39,31}, {-37,23}, {-32,12}, {-23,0}, {-12,-12}, {0,-23}, {13,-31}, {23,-37}, {31,-39}, {37,-37}, {39,-31}, {37,-23}, {31,-13}, {23,0}, {12,12}, {0,23}, {-12,32}, {-23,37}, {-31,39},
{-37,37}, {-32,31}, {-23,23}, {-12,12}, {0,0}, {12,-12}, {24,-22}, {31,-32}, {38,-37}, {40,-40}, {37,-38}, {32,-31}, {22,-24}, {12,-12}, {0,0}, {-12,12}, {-23,23}, {-31,32}, {-37,37},
{-32,39}, {-24,37}, {-12,32}, {0,23}, {12,12}, {23,0}, {32,-11}, {37,-23}, {40,-31}, {37,-37}, {31,-40}, {23,-37}, {12,-32}, {0,-24}, {-12,-12}, {-23,0}, {-31,12}, {-37,23}, {-39,32},
{-24,37}, {-12,39}, {0,37}, {12,32}, {23,23}, {32,12}, {38,0}, {39,-12}, {37,-22}, {31,-31}, {23,-37}, {12,-40}, {0,-37}, {-12,-32}, {-22,-24}, {-31,-13}, {-37,0}, {-39,12}, {-37,23},
{-12,31}, {0,37}, {12,39}, {24,37}, {32,31}, {37,23}, {39,13}, {37,0}, {32,-12}, {23,-23}, {12,-31}, {0,-37}, {-12,-40}, {-23,-37}, {-32,-31}, {-38,-23}, {-39,-12}, {-37,0}, {-32,12},
{0,24}, {12,31}, {24,37}, {32,39}, {37,37}, {39,31}, {38,23}, {32,12}, {23,0}, {12,-12}, {0,-23}, {-12,-31}, {-23,-37}, {-31,-40}, {-37,-38}, {-39,-32}, {-37,-23}, {-32,-12}, {-24,0},
{13,13}, {23,23}, {32,32}, {37,37}, {39,39}, {37,37}, {32,32}, {23,23}, {12,12}, {0,0}, {-12,-12}, {-23,-23}, {-31,-31}, {-37,-37}, {-40,-40}, {-37,-37}, {-32,-32}, {-23,-23}, {-13,-13},
{24,0}, {32,12}, {37,23}, {39,32}, {37,37}, {31,39}, {23,38}, {12,32}, {0,23}, {-12,12}, {-23,0}, {-32,-12}, {-37,-22}, {-40,-31}, {-38,-37}, {-31,-39}, {-23,-38}, {-12,-32}, {0,-24},
{32,-12}, {37,0}, {39,12}, {37,23}, {31,32}, {23,37}, {12,39}, {0,38}, {-12,32}, {-23,23}, {-32,12}, {-37,0}, {-39,-12}, {-37,-23}, {-31,-32}, {-23,-37}, {-13,-39}, {0,-37}, {12,-32},
{38,-23}, {39,-13}, {37,0}, {31,12}, {23,23}, {12,32}, {0,37}, {-12,39}, {-23,38}, {-32,32}, {-38,23}, {-39,13}, {-38,0}, {-32,-11}, {-24,-22}, {-13,-31}, {0,-37}, {13,-39}, {23,-38},
{39,-31}, {37,-23}, {31,-13}, {23,0}, {12,12}, {0,23}, {-12,32}, {-23,37}, {-31,39}, {-37,37}, {-39,31}, {-37,23}, {-32,12}, {-23,0}, {-12,-12}, {0,-23}, {13,-31}, {23,-37}, {31,-39},
{37,-38}, {32,-31}, {22,-24}, {12,-12}, {0,0}, {-12,12}, {-23,23}, {-31,32}, {-37,37}, {-39,39}, {-37,37}, {-32,31}, {-23,23}, {-12,12}, {0,0}, {12,-12}, {24,-22}, {31,-32}, {38,-37},
{31,-40}, {23,-37}, {12,-32}, {0,-24}, {-12,-12}, {-23,0}, {-31,12}, {-37,23}, {-39,32}, {-37,37}, {-32,39}, {-24,37}, {-12,32}, {0,23}, {12,12}, {23,0}, {32,-11}, {37,-23}, {40,-31},
{23,-37}, {12,-40}, {0,-37}, {-12,-32}, {-22,-24}, {-31,-13}, {-37,0}, {-39,12}, {-37,23}, {-32,32}, {-24,37}, {-12,39}, {0,37}, {12,32}, {23,23}, {32,12}, {38,0}, {39,-12}, {37,-22},
{12,-31}, {0,-37}, {-12,-40}, {-23,-37}, {-32,-31}, {-38,-23}, {-39,-12}, {-37,0}, {-32,12}, {-23,23}, {-12,31}, {0,37}, {12,39}, {24,37}, {32,31}, {37,23}, {39,13}, {37,0}, {32,-12},
{0,-23}, {-12,-31}, {-23,-37}, {-31,-40}, {-37,-38}, {-39,-32}, {-37,-23}, {-32,-12}, {-24,0}, {-13,13}, {0,24}, {12,31}, {24,37}, {32,39}, {37,37}, {39,31}, {38,23}, {32,12}, {23,0},
}

func TestFdgrad(t *testing.T) {
	var tests = []struct {
		name string
		im   SippImage
		w	 int
		h	 int
		gr   *[]complex128
		mm   float64
	} {
		{
			"Default kernel Sgray",
			Sgray,
			3,
			3,
			&smallPicGrad,
			smallPicGradMaxMod,
		},
		{
			"Default kernel SgrayCosxCosyTiny",
			SgrayCosxCosyTiny,
			19,
			19,
			&CosxCosyTinyGrad,
			CosxCosyTinyGradMaxMod,
		},
	}
	for _, test := range tests {
		grad := Fdgrad(test.im)
		if grad.Rect.Dx() != test.w ||
		   grad.Rect.Dy() != test.h {
		    t.Errorf("Error in test %s: Gradient image rect wrong size: expected [%d,%d], got [%d,%d]",
			    test.name, test.w, test.h, grad.Rect.Dx(), grad.Rect.Dy())
		}
		if !reflect.DeepEqual(grad.Pix, *test.gr) {
			t.Error("Error in test ", test.name, ": Gradient image incorrect. Expected:" +
				ComplexArrayToString(*test.gr, test.w) + "Got:" +
				ComplexArrayToString(grad.Pix, test.w))
		}
		if grad.MaxMod != test.mm {
			t.Errorf("Error in test %s: Incorrect max modulus. Expected: %f, got %f",
				test.name, test.mm, grad.MaxMod)
		}
	}
}

func TestFdgradInt32(t *testing.T) {
	var tests = []struct {
		name string
		im   SippImage
		w	 int
		h	 int
		gr   *[]ComplexInt32
	} {
		{
			"Default kernel Sgray",
			Sgray,
			3,
			3,
			&smallPicGradInt32,
		},
		{
			"Default kernel SgrayCosxCosyTiny",
			SgrayCosxCosyTiny,
			19,
			19,
			&CosxCosyTinyGradInt32,
		},
	}
	for _, test := range tests {
		grad := FdgradInt32(test.im)
		if grad.Rect.Dx() != test.w ||
		   grad.Rect.Dy() != test.h {
		    t.Errorf("Error in test %s: Gradient image rect wrong size: expected [%d,%d], got [%d,%d]",
			    test.name, test.w, test.h, grad.Rect.Dx(), grad.Rect.Dy())
		}
		if !reflect.DeepEqual(grad.Pix, *test.gr) {
			t.Error("Error in test ", test.name, ": Gradient image incorrect. Expected:" +
				ComplexInt32ArrayToString(*test.gr, test.w) + "Got:" +
				ComplexInt32ArrayToString(grad.Pix, test.w))
		}
	}
}

// equalSubImage returns true if the "test" image is a subset of either the real
// or imaginary component of the "target" complex image, in the top-left corner.
func equalSubImage(test *ComplexImage, target *SippGray, testReal bool) bool {
	for y := 0; y < test.Rect.Dy(); y++ {
		for x := 0; x < test.Rect.Dx(); x++ {
			exp := target.Val(x, y)
			tst := test.Pix[y*test.Rect.Dy()+x]
			if testReal && (real(tst) != exp) {
				return false
			}
			if (!testReal) && (imag(tst) != exp) {
				return false
			}
		}
	}
	return true
}

// equalSubImageInt32 returns true if the "test" image is a subset of either the
// real or imaginary component of the "target" complex image, in the top-left
// corner.
func equalSubImageInt32(test *ComplexInt32Image, target *SippGray, testReal bool) bool {
	for y := 0; y < test.Rect.Dy(); y++ {
		for x := 0; x < test.Rect.Dx(); x++ {
			exp := target.IntVal(x, y)
			tst := test.Pix[y*test.Rect.Dy()+x]
			if testReal && (tst.Re != exp) {
				return false
			}
			if (!testReal) && (tst.Im != exp) {
				return false
			}
		}
	}
	return true
}

func TestFdgradKernel(t *testing.T) {
	// SgrayCosxCosyTiny with identityKernel
	idGrad := FdgradKernel(SgrayCosxCosyTiny, identityKernel)
	if !equalSubImage(idGrad, SgrayCosxCosyTiny, true) {
		t.Error("Error: Real Identity gradient incorrect")
	}

	// SgrayCosxCosyTiny with imagIdentityKernel
	imagIdGrad := FdgradKernel(SgrayCosxCosyTiny, imagIdentityKernel)
	if !equalSubImage(imagIdGrad, SgrayCosxCosyTiny, false) {
		t.Error("Error: Imaginary Identity gradient incorrect")
	}

	// Sgray with kieransKernel
	kierGrad := FdgradKernel(Sgray, kieransKernel)
	if !reflect.DeepEqual(kierGrad.Pix, smallPicGradKieransKernel) {
		t.Error("Error: Gradient image incorrect. Expected:" +
			ComplexArrayToString(smallPicGradKieransKernel, 3) + "Got:" +
			ComplexArrayToString(kierGrad.Pix, 3))
	}
}

func TestFdgradKernelInt32(t *testing.T) {
	// SgrayCosxCosyTiny with identityKernel
	idGrad := FdgradInt32Kernel(SgrayCosxCosyTiny, identityKernelInt32)
	if !equalSubImageInt32(idGrad, SgrayCosxCosyTiny, true) {
		t.Error("Error: Real Identity gradient incorrect")
	}

	// SgrayCosxCosyTiny with imagIdentityKernel
	imagIdGrad := FdgradInt32Kernel(SgrayCosxCosyTiny, imagIdentityKernelInt32)
	if !equalSubImageInt32(imagIdGrad, SgrayCosxCosyTiny, false) {
		t.Error("Error: Imaginary Identity gradient incorrect")
	}

	// Sgray with kieransKernel
	kierGrad := FdgradInt32Kernel(Sgray, kieransKernelInt32)
	if !reflect.DeepEqual(kierGrad.Pix, smallPicGradKieransKernelInt32) {
		t.Error("Error: Gradient image incorrect. Expected:" +
			ComplexInt32ArrayToString(smallPicGradKieransKernelInt32, 3) + "Got:" +
			ComplexInt32ArrayToString(kierGrad.Pix, 3))
	}
}
