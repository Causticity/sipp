// Copyright Raul Vera 2020

// Tests for package shist.

package shist

import (
	_ "image/png"
	"path/filepath"
	"reflect"
	"testing"
)

import (
	. "github.com/Causticity/sipp/scomplex"
	//. "github.com/Causticity/sipp/sgrad"
	. "github.com/Causticity/sipp/simage"
	. "github.com/Causticity/sipp/sipptesting"
	. "github.com/Causticity/sipp/sipptesting/sipptestcore"
)

// TODO: The coverage tool shows a few minor code paths not tested here. Most
// are related to 16-bit, which will be reworked anyway. Once 16-bit is fully
// sorted out, the coverage should be improved to 100%.

var cosxCosyTinyBinIndex = []int{
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

var cosxCosyTinyBinIndexRadius255 = []int{
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

var expectedCosxCosyTinyNonZeroHistCount = 127

var cosxCosyTinyBinVals = []uint32{
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

var expectedMax uint32 = 8

// 0th entry is 0, then next 16 entries should have a 1, all others 0
func checkHist(t *testing.T, hist []uint32) {
	for i, val := range hist {
		var check uint32
		if i > 0 && i < 17 {
			check = 1
		} else {
			check = 0
		}
		if val != check {
			t.Errorf("Error: histogram at index %v incorrect, expected %v, got %v",
				i, check, val)
		}
	}
}

func TestGreyHist(t *testing.T) {
	hist := GreyHist(Sgray)
	checkHist(t, hist)
	hist = GreyHist(Sgray16)
	checkHist(t, hist)
}

type histTest struct {
	name                  string
	grad                  []complex128
	radius                uint16
	stride                int
	maxMod                float64
	count                 int
	binIndex              []int
	binVals               []uint32
	maxBinVal             uint32
	suppressedImageName   string
	renderedHistogramName string
}

func TestHist(t *testing.T) {
	var tests = []histTest{
		{
			"CosxCosyTinyGrad RADIUS = 0",
			CosxCosyTinyGrad,
			0,
			CosxCosyTinyStride,
			CosxCosyTinyGradMaxMod,
			expectedCosxCosyTinyNonZeroHistCount,
			cosxCosyTinyBinIndex,
			cosxCosyTinyBinVals,
			expectedMax,
			"cosxcosy_tiny_hist_sup.png",
			"cosxcosy_tiny_hist.png",
		},
		{
			"CosxCosyTinyGrad RADIUS = 255",
			CosxCosyTinyGrad,
			255,
			CosxCosyTinyStride,
			CosxCosyTinyGradMaxMod,
			expectedCosxCosyTinyNonZeroHistCount,
			cosxCosyTinyBinIndexRadius255,
			cosxCosyTinyBinVals,
			expectedMax,
			"cosxcosy_tiny_radius255_hist_sup.png",
			"cosxcosy_tiny_radius255_hist.png",
		},
		// TODO 16-bit tests, once 16-bit has been truly sorted out.
	}
	for _, test := range tests {
		grad := FromComplexArray(test.grad, test.stride-1)
		hist := Hist(grad, test.radius)
		if hist.Grad != grad {
			t.Errorf("Error: SippHist for %s has incorrect grad, expected %v, got %v",
				test.name, grad, hist.Grad)
		}
		xpctRadius := test.radius
		if test.radius == 0 {
			xpctRadius = uint16(test.maxMod) + radiusMargin
		}
		if hist.Radius != xpctRadius {
			t.Errorf("Error: radius for %s histogram incorrect. Expected %v, got %v",
				test.name, xpctRadius, hist.Radius)
		}
		count := 0
		for _, val := range hist.Bin {
			if val != 0 {
				count++
			}
		}
		xpctCnt := test.count
		if count != xpctCnt {
			t.Errorf("Error: Histogram for %s has incorrect number of non-zero entries: expected %v, got %v",
				test.name, xpctCnt, count)
		}
		if !reflect.DeepEqual(hist.BinIndex, test.binIndex) {
			t.Errorf("Error: hist.binIndex for %s incorrect, expected\n%v\n got\n%v\n",
				test.name, test.binIndex, hist.BinIndex)
		}
		for i, val := range hist.BinIndex {
			if hist.Bin[val] != test.binVals[i] {
				t.Errorf("Error: histogram value for %s incorrect, expected %v, got %v",
					test.name, test.binVals[i], hist.Bin[val])
			}
		}
		if hist.Max != test.maxBinVal {
			t.Errorf("Error: hist.Max for %s incorrect. Expected %v, got %v",
				test.name, test.maxBinVal, hist.Max)
		}
		supp := hist.RenderSuppressed()
		check, err := Read(filepath.Join(TestDir, test.suppressedImageName))
		if err != nil {
			t.Errorf("Error reading suppressed histogram check image: %v\n", test.suppressedImageName)
		}
		if !reflect.DeepEqual(supp.Pix(), check.Pix()) {
			t.Errorf("Error: suppressed histogram image incorrect. Expected %v, got%v\n",
				check.Pix(), supp.Pix())
		}
		rnd := hist.Render()
		check, err = Read(filepath.Join(TestDir, test.renderedHistogramName))
		if err != nil {
			t.Errorf("Error reading rendered histogram check image: %v\n", test.renderedHistogramName)
		}
		if !reflect.DeepEqual(rnd.Pix(), check.Pix()) {
			t.Errorf("Error: rendered histogram image incorrect. Expected %v, got%v\n",
				check.Pix(), rnd.Pix())
		}
	}
}
