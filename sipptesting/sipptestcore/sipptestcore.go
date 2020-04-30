// Copyright Raul Vera 2020

// Basic test infrastructure for testing the sipp package.

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

var SmallPicGrad = []complex128 {
    5 - 3i, 5 - 3i, 5 - 3i,
    5 - 3i, 5 - 3i, 5 - 3i,
    5 - 3i, 5 - 3i, 5 - 3i,
}

var SmallPicGradMaxMod = math.Sqrt(34.0)

var SmallPicEntropy = 4.0

var SmallPicEntropyImage = []uint8 {
    255, 255, 255, 255,
    255, 255, 255, 255,
    255, 255, 255, 255,
    255, 255, 255, 255,
}

var SmallPic16Entropy = 4.0

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

var CosxCosyTinyStride = 20

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

var CosxCosyTinyEntropy = 5.211838049805185

var CosxCosyTinyEntropyImage = []uint8 {
154, 255, 255, 154, 154, 255, 255, 255, 255, 154, 154, 255, 255, 255, 255, 154, 154, 255, 255, 154,
255, 154, 208, 154, 255, 255, 124, 208, 154, 255, 255, 154, 208, 182, 255, 255, 154, 208, 154, 255,
255, 208, 124, 255, 154, 154, 208, 154, 208, 255, 255, 208, 154, 208, 154, 154, 255,  52, 208, 255,
154, 154, 255, 124, 255, 255, 154, 208, 124, 255, 255, 124, 208, 154, 255, 255, 124, 255, 154, 154,
154, 255, 154, 255, 154, 154, 255, 154, 255, 255, 255, 255, 154, 255, 154, 154, 255, 154, 255, 154,
255, 255, 154, 255, 154, 154, 255, 154, 255, 154, 154, 255, 154, 255, 154, 154, 255, 154, 255, 255,
255, 182,  91, 154, 255, 255, 124, 255, 154, 154, 154, 154, 255,  52, 255, 255, 154, 208, 182, 255,
255,  91, 154,  91, 154, 154, 255, 124, 208, 255, 255,  91, 124, 255, 154, 154, 208, 154, 208, 255,
255, 154,  91, 182, 255, 255, 154, 208, 154, 255, 255, 154,  91, 154, 255, 255, 182, 208, 154, 255,
154, 255, 255, 255, 255, 154, 154, 255, 255, 154, 154, 255, 255, 154, 154, 255, 255, 255, 255, 154,
154, 255, 255, 255, 255, 154, 154, 255, 255, 154, 154, 255, 255, 154, 154, 255, 255, 255, 255, 154,
255, 154, 208, 182, 255, 255, 154, 208, 154, 255, 255, 154, 208, 154, 255, 255, 124, 208, 154, 255,
255, 208, 154, 208, 154, 154, 255,  52, 208, 255, 255, 208, 124, 255, 154, 154, 208, 154, 208, 255,
255, 124, 208, 154, 255, 255, 124, 255, 154, 154, 154, 154, 255, 124, 255, 255, 154, 208, 124, 255,
255, 255, 154, 255, 154, 154, 255, 154, 255, 154, 154, 255, 154, 255, 154, 154, 255, 154, 255, 255,
154, 255, 154, 255, 154, 154, 255, 154, 255, 255, 255, 255, 154, 255, 154, 154, 255, 154, 255, 154,
154, 154, 255,  52, 255, 255, 154, 208, 182, 255, 255, 182,  91, 154, 255, 255, 124, 255, 154, 154,
255,  91, 124, 255, 154, 154, 208, 154, 208, 255, 255,  91, 154,  91, 154, 154, 255, 124, 208, 255,
255, 154,  91, 154, 255, 255, 182, 208, 154, 255, 255, 154,  91, 182, 255, 255, 154, 208, 154, 255,
154, 255, 255, 154, 154, 255, 255, 255, 255, 154, 154, 255, 255, 255, 255, 154, 154, 255, 255, 154,
}

var CosxCosyTinyGradMaxMod = 56.568542494923804

var CosxCosyTinyBinIndex = []int {
8297, 6740, 5445, 4281, 3509, 3258, 3395, 4180, 5224, 6656, 8344, 9900, 11324,
12487, 13130, 13382, 13245, 12460, 11287, 9836, 8283, 6733, 5316, 4161, 3524,
3276, 3547, 4204, 5376, 6804, 8357, 9907, 11324, 12479, 13116, 13363, 13222,
12436, 11249, 9958, 8282, 6869, 5458, 4308, 3547, 3302, 3441, 4224, 5391, 6682,
8357, 9899, 11310, 12460, 13093, 13339, 13199, 12280, 11250, 9836, 8297, 6760,
5353, 4334, 3570, 3320, 3584, 4360, 5390, 6674, 8343, 9880, 11287, 12436, 13070,
13320, 13056, 12287, 11264, 9856, 8320, 6784, 5506, 4223, 3585, 3200, 3455, 4353,
5246, 6784, 8320, 9856, 11264, 12417, 13056, 13319, 13069, 12436, 11287, 9880,
8343, 6933, 5390, 4361, 3584, 3191, 3570, 4204, 5224, 6760, 8297, 9837, 11250,
12409, 13069, 13339, 13093, 12460, 11310, 9900, 8358, 6811, 5519, 4352, 3570,
3172, 3547, 4180, 5202, 6612, 8283, 9829, 11250, 12307, 13093, 13363, 13117,
12351, 11324, 10036, 8357, 6804, 5376, 4333, 3547, 3148, 3524, 4289, 5315, 6733,
8283, 9836, 11416, 12331, 13117, 13383, 13130, 12358, 11325, 9900, 8343, 6784,
5353, 4309, 3524, 3129, 3381, 4153, 5316, 6740, 8296, 10010, 11310, 12480, 13130,
13390, 13130, 12480, 11310, 9880, 8320, 6760, 5330, 4290, 3510, 3120, 3510, 4160,
5330, 6630, 8344, 9900, 11324, 12487, 13130, 13382, 13245, 12460, 11287, 9856,
8297, 6740, 5445, 4281, 3509, 3258, 3395, 4180, 5224, 6804, 8357, 9907, 11324,
12479, 13116, 13363, 13222, 12436, 11264, 9836, 8283, 6733, 5316, 4161, 3524,
3276, 3547, 4204, 5391, 6682, 8357, 9899, 11310, 12460, 13093, 13339, 13199,
12416, 11249, 9958, 8282, 6869, 5458, 4308, 3547, 3302, 3441, 4360, 5390, 6674,
8343, 9880, 11287, 12436, 13070, 13320, 13056, 12280, 11250, 9836, 8297, 6760,
5353, 4334, 3570, 3320, 3455, 4353, 5246, 6784, 8320, 9856, 11264, 12417, 13056,
13312, 13056, 12287, 11264, 9856, 8320, 6784, 5506, 4223, 3585, 3191, 3570, 4204,
5224, 6760, 8297, 9837, 11250, 12409, 13056, 13319, 13069, 12436, 11287, 9880,
8343, 6933, 5390, 4361, 3570, 3172, 3547, 4180, 5202, 6612, 8283, 9829, 11250,
12416, 13069, 13339, 13093, 12460, 11310, 9900, 8358, 6811, 5519, 4333, 3547,
3148, 3524, 4289, 5315, 6733, 8283, 9836, 11264, 12307, 13093, 13363, 13117,
12351, 11324, 10036, 8357, 6804, 5353, 4309, 3524, 3129, 3381, 4153, 5316, 6740,
8296, 9984, 11416, 12331, 13117, 13383, 13130, 12358, 11325, 9900, 8343,
}

var CosxCosyTinyBinIndexk255 = []int {
130537, 124396, 119281, 114679, 111615, 110600, 111119, 114196, 118296, 123930,
130584, 136724, 142350, 146951, 149504, 150520, 150001, 146924, 142313, 136660,
130523, 124389, 118770, 114177, 111630, 110618, 111653, 114220, 118830, 124460,
130597, 136731, 142350, 146943, 149490, 150501, 149978, 146900, 142275, 137164,
130522, 124907, 119294, 114706, 111653, 110644, 111165, 114240, 118845, 123956,
130597, 136723, 142336, 146924, 149467, 150477, 149955, 146362, 142276, 136660,
130537, 124416, 118807, 114732, 111676, 110662, 111690, 114758, 118844, 123948,
130583, 136704, 142313, 146900, 149444, 150458, 149430, 146369, 142290, 136680,
130560, 124440, 119342, 114239, 111691, 110160, 111179, 114751, 118318, 124440,
130560, 136680, 142290, 146881, 149430, 150457, 149443, 146900, 142313, 136704,
130583, 124971, 118844, 114759, 111690, 110151, 111676, 114220, 118296, 124416,
130537, 136661, 142276, 146873, 149443, 150477, 149467, 146924, 142336, 136724,
130598, 124467, 119355, 114750, 111676, 110132, 111653, 114196, 118274, 123886,
130523, 136653, 142276, 146389, 149467, 150501, 149491, 146433, 142350, 137242,
130597, 124460, 118830, 114731, 111653, 110108, 111630, 114687, 118769, 124389,
130523, 136660, 142824, 146413, 149491, 150521, 149504, 146440, 142351, 136724,
130583, 124440, 118807, 114707, 111630, 110089, 111105, 114169, 118770, 124396,
130536, 137216, 142336, 146944, 149504, 150528, 149504, 146944, 142336, 136704,
130560, 124416, 118784, 114688, 111616, 110080, 111616, 114176, 118784, 123904,
130584, 136724, 142350, 146951, 149504, 150520, 150001, 146924, 142313, 136680,
130537, 124396, 119281, 114679, 111615, 110600, 111119, 114196, 118296, 124460,
130597, 136731, 142350, 146943, 149490, 150501, 149978, 146900, 142290, 136660,
130523, 124389, 118770, 114177, 111630, 110618, 111653, 114220, 118845, 123956,
130597, 136723, 142336, 146924, 149467, 150477, 149955, 146880, 142275, 137164,
130522, 124907, 119294, 114706, 111653, 110644, 111165, 114758, 118844, 123948,
130583, 136704, 142313, 146900, 149444, 150458, 149430, 146362, 142276, 136660,
130537, 124416, 118807, 114732, 111676, 110662, 111179, 114751, 118318, 124440,
130560, 136680, 142290, 146881, 149430, 150450, 149430, 146369, 142290, 136680,
130560, 124440, 119342, 114239, 111691, 110151, 111676, 114220, 118296, 124416,
130537, 136661, 142276, 146873, 149430, 150457, 149443, 146900, 142313, 136704,
130583, 124971, 118844, 114759, 111676, 110132, 111653, 114196, 118274, 123886,
130523, 136653, 142276, 146880, 149443, 150477, 149467, 146924, 142336, 136724,
130598, 124467, 119355, 114731, 111653, 110108, 111630, 114687, 118769, 124389,
130523, 136660, 142290, 146389, 149467, 150501, 149491, 146433, 142350, 137242,
130597, 124460, 118807, 114707, 111630, 110089, 111105, 114169, 118770, 124396,
130536, 137190, 142824, 146413, 149491, 150521, 149504, 146440, 142351, 136724,
130583,
}

var ExpectedCosxCosyTinyNonZeroHistCount = 127

var CosxCosyTinyBinVals = []uint32 {
6, 4, 2, 2, 2, 2, 2, 4, 4, 1, 2, 6, 6, 2, 6, 2, 2, 6, 6,
6, 6, 4, 4, 2, 6, 2, 8, 4, 2, 4, 6, 2, 6, 2, 2, 4, 2, 6,
2, 2, 2, 2, 2, 2, 8, 2, 2, 1, 2, 2, 6, 2, 6, 6, 6, 4, 2,
2, 6, 6, 6, 5, 4, 2, 6, 2, 2, 2, 4, 2, 6, 5, 6, 6, 2, 2,
6, 2, 6, 5, 5, 5, 2, 2, 2, 1, 2, 2, 2, 5, 5, 5, 6, 2, 6,
2, 4, 6, 6, 5, 6, 2, 4, 2, 2, 2, 6, 4, 4, 5, 6, 2, 6, 2,
4, 4, 6, 6, 6, 6, 2, 2, 2, 1, 6, 2, 8, 4, 2, 2, 6, 2, 6,
2, 6, 4, 4, 2, 6, 2, 6, 4, 2, 2, 8, 2, 6, 2, 2, 4, 6, 6,
2, 2, 4, 2, 6, 2, 2, 6, 6, 5, 4, 2, 6, 2, 2, 2, 4, 4, 2,
1, 6, 2, 6, 1, 6, 2, 6, 5, 5, 5, 2, 1, 2, 1, 2, 1, 2, 1,
2, 6, 6, 2, 6, 2, 2, 6, 6, 5, 6, 4, 2, 2, 2, 2, 2, 4, 4,
4, 6, 2, 6, 2, 2, 4, 2, 6, 6, 6, 6, 4, 4, 2, 6, 2, 8, 4,
2, 2, 6, 2, 6, 6, 6, 4, 2, 2, 2, 2, 2, 2, 2, 2, 8, 2, 2,
2, 4, 2, 6, 5, 6, 6, 2, 2, 6, 2, 6, 6, 6, 5, 4, 2, 6, 2,
2, 2, 2, 5, 5, 5, 6, 2, 6, 1, 6, 2, 6, 5, 5, 5, 2, 2, 2,
2, 6, 4, 4, 5, 6, 2, 6, 2, 6, 2, 4, 6, 6, 5, 6, 2, 4, 2,
6, 2, 8, 4, 2, 2, 6, 2, 6, 2, 4, 4, 6, 6, 6, 6, 2, 2, 2,
2, 8, 2, 6, 2, 2, 4, 6, 6, 6, 2, 6, 4, 4, 2, 6, 2, 6, 4,
4, 2, 6, 2, 2, 2, 4, 4, 2, 1, 2, 2, 4, 2, 6, 2, 2, 6, 6,
}

var ExpectedMax uint32 = 8
var ExpectedDelentropy float64 = 6.775012499324645
var ExpectedMaxDelentropy float64 = 0.12179180114985422
var ExpectedDelentropyArray = []float64 {
0, 0.023534224451211, 0.04152828269743585, 0, 0.0719762329848994,
0.0855114533517979, 0.09824198104431049, 0, 0.12179180114985422,
}

var cosxCosyTinyDelentropyImageArray = []uint8 {
205, 150,  86,  86,  86,  86,  86, 150, 150,  49,  86, 205, 205,  86, 205,  86,  86, 205, 205,
205, 205, 150, 150,  86, 205,  86, 255, 150,  86, 150, 205,  86, 205,  86,  86, 150,  86, 205,
 86,  86,  86,  86,  86,  86, 255,  86,  86,  49,  86,  86, 205,  86, 205, 205, 205, 150,  86,
 86, 205, 205, 205, 179, 150,  86, 205,  86,  86,  86, 150,  86, 205, 179, 205, 205,  86,  86,
205,  86, 205, 179, 179, 179,  86,  86,  86,  49,  86,  86,  86, 179, 179, 179, 205,  86, 205,
 86, 150, 205, 205, 179, 205,  86, 150,  86,  86,  86, 205, 150, 150, 179, 205,  86, 205,  86,
150, 150, 205, 205, 205, 205,  86,  86,  86,  49, 205,  86, 255, 150,  86,  86, 205,  86, 205,
 86, 205, 150, 150,  86, 205,  86, 205, 150,  86,  86, 255,  86, 205,  86,  86, 150, 205, 205,
 86,  86, 150,  86, 205,  86,  86, 205, 205, 179, 150,  86, 205,  86,  86,  86, 150, 150,  86,
 49, 205,  86, 205,  49, 205,  86, 205, 179, 179, 179,  86,  49,  86,  49,  86,  49,  86,  49,
 86, 205, 205,  86, 205,  86,  86, 205, 205, 179, 205, 150,  86,  86,  86,  86,  86, 150, 150,
150, 205,  86, 205,  86,  86, 150,  86, 205, 205, 205, 205, 150, 150,  86, 205,  86, 255, 150,
 86,  86, 205,  86, 205, 205, 205, 150,  86,  86,  86,  86,  86,  86,  86,  86, 255,  86,  86,
 86, 150,  86, 205, 179, 205, 205,  86,  86, 205,  86, 205, 205, 205, 179, 150,  86, 205,  86,
 86,  86,  86, 179, 179, 179, 205,  86, 205,  49, 205,  86, 205, 179, 179, 179,  86,  86,  86,
 86, 205, 150, 150, 179, 205,  86, 205,  86, 205,  86, 150, 205, 205, 179, 205,  86, 150,  86,
205,  86, 255, 150,  86,  86, 205,  86, 205,  86, 150, 150, 205, 205, 205, 205,  86,  86,  86,
 86, 255,  86, 205,  86,  86, 150, 205, 205, 205,  86, 205, 150, 150,  86, 205,  86, 205, 150,
150,  86, 205,  86,  86,  86, 150, 150,  86,  49,  86,  86, 150,  86, 205,  86,  86, 205, 205,
}

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

var GrayZero = image.Gray {
    smallZeroPic, 
    4, 
    image.Rectangle{image.Point{0, 0}, image.Point{4, 4}},
}

var CosxCosyTiny = image.Gray {
    cosxCosyTiny,
    20,
    image.Rectangle{image.Point{0, 0}, image.Point{20,20}},
}

var CosxCosyTinyDelentropyImage = image.Gray {
    cosxCosyTinyDelentropyImageArray,
    19,
    image.Rectangle{image.Point{0, 0}, image.Point{19,19}},
}