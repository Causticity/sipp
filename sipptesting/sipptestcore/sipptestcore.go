// Copyright Raul Vera 2020

// Tests for package sgrad.

package sipptestcore

import (
    "fmt"
	"image"
	"math"
	"os"
	"path/filepath"
)

var TestDir = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", 
							"Causticity", "sipp", "testdata")

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

var ShiftedPic = []complex128 {
    1,  -2,  3,  -4,
    -5,  6,  -7,  8,
    9, -10, 11, -12,
    -13, 14, -15, 16,
}

var ScaledShiftedPic = []uint8 {
    131, 106, 148, 90,
    82, 172, 65, 189, 
    197, 41, 213, 24,
    16, 238, 0, 255,
}

// TODO: This gradient is too simple. Test with a more complex image
var SmallPicGrad = []complex128 {
    5 - 3i, 5 - 3i, 5 - 3i,
    5 - 3i, 5 - 3i, 5 - 3i,
    5 - 3i, 5 - 3i, 5 - 3i,
}

func ComplexArrayToString(cpx []complex128, stride int) string {
    res := "\n"
    rows := len(cpx)/stride
    for y := 0; y < rows; y++ {
        for x := 0; x < stride; x++ {
            res += fmt.Sprintf(" %v,", cpx[y*stride + x])
        }
        res += "\n"
    }
    return res
}

var SmallPicGradMaxMod = math.Sqrt(34.0) 

var Gray = image.Gray {
    smallPic, 
    4, 
    image.Rectangle{image.Point{0, 0}, image.Point{4, 4}},
}

var Gray16 = image.Gray16 {
    smallPic16, 
    8, 
    image.Rectangle{image.Point{0, 0}, image.Point{4, 4}},
}

var GrayZero = image.Gray {
    smallZeroPic, 
    4, 
    image.Rectangle{image.Point{0, 0}, image.Point{4, 4}},
}

