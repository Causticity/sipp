// Copyright Raul Vera 2016

// Tests for package scomplex

package scomplex

import (
	. "image"
	"reflect"
	"testing"
	)

import (
	. "github.com/Causticity/sipp/simage"
)

var smallPic = []uint8 {
    1,  2,  3,  4,
    5,  6,  7,  8,
    9, 10, 11, 12,
    13, 14, 15, 16,
}
 
var smallPic16 = []uint8 {
    0, 1, 0, 2, 0, 3, 0, 4,
    0, 5, 0, 6, 0, 7, 0, 8,
    0, 9, 0, 10, 0, 11, 0, 12,
    0, 13, 0, 14, 0, 15, 0, 16,
}
 
var smallZeroPic = []uint8 {
    0,  0,  0,  0,
    0,  0,  0,  0,
    0,  0,  0,  0,
    0,  0,  0,  0,
}

var shiftedPic = []complex128 {
    1,  -2,  3,  -4,
    -5,  6,  -7,  8,
    9, -10, 11, -12,
    -13, 14, -15, 16,
}

var scaledShiftedPic = []uint8 {
    131, 106, 148, 90,
    82, 172, 65, 189, 
    197, 41, 213, 24,
    16, 238, 0, 255,
}
 
var gray = Gray {
    smallPic, 
    4, 
    Rectangle{Point{0, 0}, Point{4, 4}},
}

var gray16 = Gray16 {
    smallPic16, 
    8, 
    Rectangle{Point{0, 0}, Point{4, 4}},
}

var grayZero = Gray {
    smallZeroPic, 
    4, 
    Rectangle{Point{0, 0}, Point{4, 4}},
}

var sgray *SippGray
var sgray16 *SippGray16
var sgrayZero *SippGray

func init() {
    sgray = new(SippGray)
    sgray.Gray = &gray
    sgray16 = new(SippGray16)
    sgray16.Gray16 = &gray16
    sgrayZero = new(SippGray)
    sgrayZero.Gray = &grayZero
}

func TestComplex (t *testing.T) {
    comp := ToShiftedComplex(sgray)
    if !reflect.DeepEqual(shiftedPic, comp.Pix) {
        t.Error("Shifted complex doesn't match Gray!")
    }
    
    comp = ToShiftedComplex(sgray16)
    if !reflect.DeepEqual(shiftedPic, comp.Pix) {
        t.Error("Shifted complex doesn't match Gray16!")
    }

    re, im := comp.Render()
    
    if !reflect.DeepEqual(re.Pix(), scaledShiftedPic) {
        t.Error("real unexpected")
    }
    
    if !reflect.DeepEqual(im, sgrayZero) {
        t.Error("imaginary not zero")
    }
}

