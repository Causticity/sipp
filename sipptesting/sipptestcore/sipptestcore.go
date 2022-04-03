// Copyright Raul Vera 2020

// Basic test infrastructure for testing the sipp package.

package sipptestcore

import (
    "fmt"
    "image"
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
 
var cosxCosyTiny = []uint8 {
252, 240, 217, 186, 149, 109,  72,  40,  17,   4,   4,  17,  40,  72, 109, 148, 185, 217, 240, 252,
240, 229, 208, 180, 146, 111,  78,  49,  28,  17,  17,  28,  49,  77, 111, 146, 179, 208, 229, 240,
217, 208, 192, 169, 143, 115,  88,  65,  49,  40,  40,  49,  65,  88, 114, 142, 169, 191, 208, 217,
185, 179, 169, 154, 137, 119, 102,  88,  78,  72,  72,  78,  88, 102, 119, 137, 154, 169, 179, 185,
148, 146, 142, 137, 131, 125, 119, 115, 111, 109, 109, 111, 115, 119, 125, 131, 137, 142, 146, 148,
109, 111, 114, 119, 125, 131, 137, 143, 146, 149, 149, 146, 143, 137, 131, 125, 119, 114, 111, 109,
 72,  77,  87, 102, 119, 137, 154, 169, 180, 186, 186, 180, 169, 155, 137, 119, 102,  88,  77,  72,
 40,  48,  65,  87, 114, 142, 169, 192, 208, 217, 217, 209, 192, 169, 143, 115,  88,  65,  49,  40,
 17,  28,  48,  77, 111, 146, 179, 208, 229, 240, 240, 229, 209, 180, 146, 111,  77,  49,  28,  17,
  4,  17,  40,  72, 109, 148, 185, 217, 240, 252, 252, 240, 217, 186, 149, 109,  72,  40,  17,   4,
  4,  17,  40,  72, 109, 148, 185, 217, 240, 252, 252, 240, 217, 186, 149, 109,  72,  40,  17,   4,
 17,  28,  49,  77, 111, 146, 179, 208, 229, 240, 240, 229, 208, 180, 146, 111,  78,  49,  28,  17,
 40,  49,  65,  88, 114, 142, 169, 191, 208, 217, 217, 208, 192, 169, 143, 115,  88,  65,  49,  40,
 72,  78,  88, 102, 119, 137, 154, 169, 179, 185, 185, 179, 169, 154, 137, 119, 102,  88,  78,  72,
109, 111, 115, 119, 125, 131, 137, 142, 146, 148, 148, 146, 142, 137, 131, 125, 119, 115, 111, 109,
149, 146, 143, 137, 131, 125, 119, 114, 111, 109, 109, 111, 114, 119, 125, 131, 137, 143, 146, 149,
186, 180, 169, 155, 137, 119, 102,  88,  77,  72,  72,  77,  87, 102, 119, 137, 154, 169, 180, 186,
217, 209, 192, 169, 143, 115,  88,  65,  49,  40,  40,  48,  65,  87, 114, 142, 169, 192, 208, 217,
240, 229, 209, 180, 146, 111,  77,  49,  28,  17,  17,  28,  48,  77, 111, 146, 179, 208, 229, 240,
252, 240, 217, 186, 149, 109,  72,  40,  17,   4,   4,  17,  40,  72, 109, 148, 185, 217, 240, 252,
}

var CosxCosyTinyGrad = []complex128 {
(-23+0i), (-32-12i), (-37-22i), (-40-31i), (-38-37i), (-31-39i), (-23-38i), (-12-32i), (0-24i), (13-13i), (24+0i), (32+12i), (37+23i), (39+32i), (37+37i), (31+39i), (23+38i), (12+32i), (0+23i),
(-32+12i), (-37+0i), (-39-12i), (-37-23i), (-31-32i), (-23-37i), (-13-39i), (0-37i), (12-32i), (23-23i), (32-12i), (37+0i), (39+12i), (37+23i), (31+32i), (23+37i), (12+39i), (0+38i), (-12+32i),
(-38+23i), (-39+13i), (-38+0i), (-32-11i), (-24-22i), (-13-31i), (0-37i), (13-39i), (23-38i), (32-32i), (38-23i), (39-13i), (37+0i), (31+12i), (23+23i), (12+32i), (0+37i), (-12+39i), (-23+38i),
(-39+31i), (-37+23i), (-32+12i), (-23+0i), (-12-12i), (0-23i), (13-31i), (23-37i), (31-39i), (37-37i), (39-31i), (37-23i), (31-13i), (23+0i), (12+12i), (0+23i), (-12+32i), (-23+37i), (-31+39i),
(-37+37i), (-32+31i), (-23+23i), (-12+12i), (0+0i), (12-12i), (24-22i), (31-32i), (38-37i), (40-40i), (37-38i), (32-31i), (22-24i), (12-12i), (0+0i), (-12+12i), (-23+23i), (-31+32i), (-37+37i),
(-32+39i), (-24+37i), (-12+32i), (0+23i), (12+12i), (23+0i), (32-11i), (37-23i), (40-31i), (37-37i), (31-40i), (23-37i), (12-32i), (0-24i), (-12-12i), (-23+0i), (-31+12i), (-37+23i), (-39+32i),
(-24+37i), (-12+39i), (0+37i), (12+32i), (23+23i), (32+12i), (38+0i), (39-12i), (37-22i), (31-31i), (23-37i), (12-40i), (0-37i), (-12-32i), (-22-24i), (-31-13i), (-37+0i), (-39+12i), (-37+23i),
(-12+31i), (0+37i), (12+39i), (24+37i), (32+31i), (37+23i), (39+13i), (37+0i), (32-12i), (23-23i), (12-31i), (0-37i), (-12-40i), (-23-37i), (-32-31i), (-38-23i), (-39-12i), (-37+0i), (-32+12i),
(0+24i), (12+31i), (24+37i), (32+39i), (37+37i), (39+31i), (38+23i), (32+12i), (23+0i), (12-12i), (0-23i), (-12-31i), (-23-37i), (-31-40i), (-37-38i), (-39-32i), (-37-23i), (-32-12i), (-24+0i),
(13+13i), (23+23i), (32+32i), (37+37i), (39+39i), (37+37i), (32+32i), (23+23i), (12+12i), (0+0i), (-12-12i), (-23-23i), (-31-31i), (-37-37i), (-40-40i), (-37-37i), (-32-32i), (-23-23i), (-13-13i),
(24+0i), (32+12i), (37+23i), (39+32i), (37+37i), (31+39i), (23+38i), (12+32i), (0+23i), (-12+12i), (-23+0i), (-32-12i), (-37-22i), (-40-31i), (-38-37i), (-31-39i), (-23-38i), (-12-32i), (0-24i),
(32-12i), (37+0i), (39+12i), (37+23i), (31+32i), (23+37i), (12+39i), (0+38i), (-12+32i), (-23+23i), (-32+12i), (-37+0i), (-39-12i), (-37-23i), (-31-32i), (-23-37i), (-13-39i), (0-37i), (12-32i),
(38-23i), (39-13i), (37+0i), (31+12i), (23+23i), (12+32i), (0+37i), (-12+39i), (-23+38i), (-32+32i), (-38+23i), (-39+13i), (-38+0i), (-32-11i), (-24-22i), (-13-31i), (0-37i), (13-39i), (23-38i),
(39-31i), (37-23i), (31-13i), (23+0i), (12+12i), (0+23i), (-12+32i), (-23+37i), (-31+39i), (-37+37i), (-39+31i), (-37+23i), (-32+12i), (-23+0i), (-12-12i), (0-23i), (13-31i), (23-37i), (31-39i),
(37-38i), (32-31i), (22-24i), (12-12i), (0+0i), (-12+12i), (-23+23i), (-31+32i), (-37+37i), (-39+39i), (-37+37i), (-32+31i), (-23+23i), (-12+12i), (0+0i), (12-12i), (24-22i), (31-32i), (38-37i),
(31-40i), (23-37i), (12-32i), (0-24i), (-12-12i), (-23+0i), (-31+12i), (-37+23i), (-39+32i), (-37+37i), (-32+39i), (-24+37i), (-12+32i), (0+23i), (12+12i), (23+0i), (32-11i), (37-23i), (40-31i),
(23-37i), (12-40i), (0-37i), (-12-32i), (-22-24i), (-31-13i), (-37+0i), (-39+12i), (-37+23i), (-32+32i), (-24+37i), (-12+39i), (0+37i), (12+32i), (23+23i), (32+12i), (38+0i), (39-12i), (37-22i),
(12-31i), (0-37i), (-12-40i), (-23-37i), (-32-31i), (-38-23i), (-39-12i), (-37+0i), (-32+12i), (-23+23i), (-12+31i), (0+37i), (12+39i), (24+37i), (32+31i), (37+23i), (39+13i), (37+0i), (32-12i),
(0-23i), (-12-31i), (-23-37i), (-31-40i), (-37-38i), (-39-32i), (-37-23i), (-32-12i), (-24+0i), (-13+13i), (0+24i), (12+31i), (24+37i), (32+39i), (37+37i), (39+31i), (38+23i), (32+12i), (23+0i),
}

var CosxCosyTinyGradMaxMod = 56.568542494923804
var CosxCosyTinyGradMaxRe float64 = 40
var CosxCoxyTinyGradMinRe float64 = -40
var CosxCosyTinyGradMaxIm float64 = 39
var CosxCosyTinyGradMinIm float64 = -40

func ComplexArrayToString(cpx []complex128, width int) string {
    res := "\n"
    rows := len(cpx)/width
    for y := 0; y < rows; y++ {
        for x := 0; x < width; x++ {
            res += fmt.Sprintf(" %v,", cpx[y*width + x])
        }
        res += "\n"
    }
    return res
}

func GrayArrayToString(gray []uint8, width int) string {
    res := "\n"
    rows := len(gray)/width
    for y := 0; y < rows; y++ {
        for x := 0; x < width; x++ {
            res += fmt.Sprintf(" %v,", gray[y*width + x])
        }
        res += "\n"
    }
    return res
}

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

var CosxCosyTiny = image.Gray {
    cosxCosyTiny,
    20,
    image.Rectangle{image.Point{0, 0}, image.Point{20,20}},
}

var CosxCosyTinyStride = 20
var CosxCosyTinyMaxReExc = 40
var CosxCosyTinyMaxImExc = 40
var CosxCosyTinyMaxExcursion = 40
