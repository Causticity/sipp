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

// A small image

var smallPic = []uint8 {
    1,  2,  3,  4,
    5,  6,  7,  8,
    9, 10, 11, 12,
    13, 14, 15, 16,
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

var grayZero = Gray {
    smallZeroPic, 
    4, 
    Rectangle{Point{0, 0}, Point{4, 4}},
}

var sgray *SippGray
var sgrayZero *SippGray

func init() {
    sgray = new(SippGray)
    sgray.Gray = &gray
    sgrayZero = new(SippGray)
    sgrayZero.Gray = &grayZero
}

func TestComplex (t *testing.T) {
    comp := ToShiftedComplex(sgray)
    if !reflect.DeepEqual(shiftedPic, comp.Pix) {
        t.Error("complex doesn't match!")
    }

    re, im := comp.Render()
    
    if !reflect.DeepEqual(re.Pix(), scaledShiftedPic) {
        t.Error("real unexpected")
    }
    
    if !reflect.DeepEqual(im, sgrayZero) {
        t.Error("imaginary not zero")
    }
}

