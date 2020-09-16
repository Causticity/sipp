// Copyright Raul Vera 2016

// Tests for the ComplexInt32 components of package scomplex

package scomplex

import (
	"fmt"
	"image"
	"reflect"
	"testing"
)

import (
	. "github.com/Causticity/sipp/simage"
	. "github.com/Causticity/sipp/sipptesting"
	. "github.com/Causticity/sipp/sipptesting/sipptestcore"
)

var shiftedPicInt32 = []ComplexInt32{
	{1,0}, {-2,0}, {3,0}, {-4,0},
	{-5,0}, {6,0}, {-7,0}, {8,0},
	{9,0}, {-10,0}, {11,0}, {-12,0},
	{-13,0}, {14,0}, {-15,0}, {16,0},
}

var shiftedPicInt32MinRe int32 = -15
var shiftedPicInt32MaxRe int32 = 16
var shiftedPicInt32MaxMod float64 = 16.0

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

var CosxCosyTinyGradInt32MinRe int32 = -40
var CosxCosyTinyGradInt32MaxRe int32 = 40
var CosxCosyTinyGradInt32MinIm int32 = -40
var CosxCosyTinyGradInt32MaxIm int32 = 39

func TestFromComplexInt32(t *testing.T) {
	cpx := FromComplexInt32Array(CosxCosyTinyGradInt32, 19)
	if !reflect.DeepEqual(cpx.Pix, CosxCosyTinyGradInt32) {
		t.Error("Error: Complex image array differs from the one constructed from")
	}
	rect := image.Rect(0, 0, 19, 19)
	if !reflect.DeepEqual(cpx.Rect, rect) {
		t.Errorf("Error: Complex image rect incorrect, expected %v, got %v\n",
			rect, cpx.Rect)
	}
	if cpx.MinRe != CosxCosyTinyGradInt32MinRe {
		t.Errorf("Error: Incorrect minimum real value. Expected: %v, got %v",
			CosxCosyTinyGradInt32MinRe, cpx.MinRe)
	}
	if cpx.MaxRe != CosxCosyTinyGradInt32MaxRe {
		t.Errorf("Error: Incorrect maximum real value. Expected: %v, got %v",
			CosxCosyTinyGradInt32MaxRe, cpx.MaxRe)
	}
	if cpx.MinIm != CosxCosyTinyGradInt32MinIm {
		t.Errorf("Error: Incorrect minimum imaginary value. Expected: %v, got %v",
			CosxCosyTinyGradInt32MinIm, cpx.MinIm)
	}
	if cpx.MaxIm != CosxCosyTinyGradInt32MaxIm {
		t.Errorf("Error: Incorrect maximum imaginary value. Expected: %v, got %v",
			CosxCosyTinyGradInt32MaxIm, cpx.MaxIm)
	}
	if cpx.MaxMod != CosxCosyTinyGradMaxMod {
		t.Errorf("Error: Incorrect max modulus. Expected: %v, got %v",
			CosxCosyTinyGradMaxMod, cpx.MaxMod)
	}
}

func TestComplexInt32Image(t *testing.T) {
	SgrayZero = new(SippGray)
	SgrayZero.Gray = &GrayZero

	comp := ToShiftedComplexInt32(Sgray)
	if !reflect.DeepEqual(shiftedPicInt32, comp.Pix) {
		t.Error("Shifted complexInt32 image doesn't match Gray!")
	}
	if comp.MinRe != shiftedPicInt32MinRe {
		t.Errorf("Shifted complexInt32 MinRe incorrect. Expected: %v, got %v",
			shiftedPicInt32MinRe, comp.MinRe)
	}
	if comp.MaxRe != shiftedPicInt32MaxRe {
		t.Errorf("Shifted complexInt32 MaxRe incorrect. Expected: %v, got %v",
			shiftedPicInt32MaxRe, comp.MaxRe)
	}
	if comp.MaxMod != shiftedPicInt32MaxMod {
		t.Errorf("Shifted complexInt32 MaxMod incorrect. Expected: %v, got %v",
			shiftedPicInt32MaxMod, comp.MaxMod)
	}

	comp = ToShiftedComplexInt32(Sgray16)
	if !reflect.DeepEqual(shiftedPicInt32, comp.Pix) {
		t.Error("Shifted complexInt32 image doesn't match Gray16!")
	}

	re, im := comp.Render()

	if !reflect.DeepEqual(re.Pix(), scaledShiftedPic) {
		fmt.Println("Real image:")
		fmt.Println(GrayArrayToString(re.Pix(), 4))
		t.Error("real unexpected")
	}

	if !reflect.DeepEqual(im, SgrayZero) {
		t.Error("imaginary not zero")
	}
}
