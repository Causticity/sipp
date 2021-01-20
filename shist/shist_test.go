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
var cosxCosyTinyBinIndex = []int{
	3257, 2276, 1461, 729, 245, 90, 179, 676, 1336, 2240, 3304, 4284, 5180,
	5911, 6314, 6470, 6381, 5884, 5143, 4220, 3243, 2269, 1380, 657, 260, 108,
	283, 700, 1440, 2340, 3317, 4291, 5180, 5903, 6300, 6451, 6358, 5860, 5105,
	4294, 3242, 2357, 1474, 756, 283, 134, 225, 720, 1455, 2266, 3317, 4283,
	5166, 5884, 6277, 6427, 6335, 5752, 5106, 4220, 3257, 2296, 1417, 782, 306,
	152, 320, 808, 1454, 2258, 3303, 4264, 5143, 5860, 6254, 6408, 6240, 5759,
	5120, 4240, 3280, 2320, 1522, 719, 321, 80, 239, 801, 1358, 2320, 3280,
	4240, 5120, 5841, 6240, 6407, 6253, 5860, 5143, 4264, 3303, 2421, 1454, 809,
	320, 71, 306, 700, 1336, 2296, 3257, 4221, 5106, 5833, 6253, 6427, 6277,
	5884, 5166, 4284, 3318, 2347, 1535, 800, 306, 52, 283, 676, 1314, 2196,
	3243, 4213, 5106, 5779, 6277, 6451, 6301, 5823, 5180, 4372, 3317, 2340, 1440,
	781, 283, 28, 260, 737, 1379, 2269, 3243, 4220, 5224, 5803, 6301, 6471,
	6314, 5830, 5181, 4284, 3303, 2320, 1417, 757, 260, 9, 165, 649, 1380,
	2276, 3256, 4346, 5166, 5904, 6314, 6478, 6314, 5904, 5166, 4264, 3280, 2296,
	1394, 738, 246, 0, 246, 656, 1394, 2214, 3304, 4284, 5180, 5911, 6314,
	6470, 6381, 5884, 5143, 4240, 3257, 2276, 1461, 729, 245, 90, 179, 676,
	1336, 2340, 3317, 4291, 5180, 5903, 6300, 6451, 6358, 5860, 5120, 4220, 3243,
	2269, 1380, 657, 260, 108, 283, 700, 1455, 2266, 3317, 4283, 5166, 5884,
	6277, 6427, 6335, 5840, 5105, 4294, 3242, 2357, 1474, 756, 283, 134, 225,
	808, 1454, 2258, 3303, 4264, 5143, 5860, 6254, 6408, 6240, 5752, 5106, 4220,
	3257, 2296, 1417, 782, 306, 152, 239, 801, 1358, 2320, 3280, 4240, 5120,
	5841, 6240, 6400, 6240, 5759, 5120, 4240, 3280, 2320, 1522, 719, 321, 71,
	306, 700, 1336, 2296, 3257, 4221, 5106, 5833, 6240, 6407, 6253, 5860, 5143,
	4264, 3303, 2421, 1454, 809, 306, 52, 283, 676, 1314, 2196, 3243, 4213,
	5106, 5840, 6253, 6427, 6277, 5884, 5166, 4284, 3318, 2347, 1535, 781, 283,
	28, 260, 737, 1379, 2269, 3243, 4220, 5120, 5779, 6277, 6451, 6301, 5823,
	5180, 4372, 3317, 2340, 1417, 757, 260, 9, 165, 649, 1380, 2276, 3256,
	4320, 5224, 5803, 6301, 6471, 6314, 5830, 5181, 4284, 3303,
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
			"CosxCosyTinyGrad",
			CosxCosyTinyGrad,
			CosxCosyTinyStride,
			CosxCosyTinyGradMaxMod,
			expectedCosxCosyTinyNonZeroHistCount,
			cosxCosyTinyBinIndex,
			cosxCosyTinyBinVals,
			expectedMax,
			"cosxcosy_tiny_hist_sup.png",
			"cosxcosy_tiny_hist.png",
		},
		// TODO 16-bit tests, once 16-bit has been truly sorted out.
	}
	for _, test := range tests {
		grad := FromComplexArray(test.grad, test.stride-1)
		hist := Hist(grad)
		if hist.Grad != grad {
			t.Errorf("Error: SippHist for %s has incorrect grad, expected %v, got %v",
				test.name, grad, hist.Grad)
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
		checkName := filepath.Join(TestDir, test.suppressedImageName)
		check, err := Read(checkName)
		if err != nil {
			t.Errorf("Error reading suppressed histogram check image: %v\n", test.suppressedImageName)
		}
		if !reflect.DeepEqual(supp.Pix(), check.Pix()) {
			// Write out the check image and report names of mismatched files
			name := SaveFailedSimage(checkName, supp)
			t.Errorf("Error: suppressed histogram image does not match expected.\nExpected in file " +
				checkName + "\nFailed saved in file " + name)
		}
		rnd := hist.Render()
		checkName = filepath.Join(TestDir, test.renderedHistogramName)
		check, err = Read(checkName)
		if err != nil {
			t.Errorf("Error reading rendered histogram check image: %v\n", test.renderedHistogramName)
		}
		if !reflect.DeepEqual(rnd.Pix(), check.Pix()) {
			name := SaveFailedSimage(checkName, rnd)
			t.Errorf("Error: rendered histogram image does not match expected.\nExpected in file " +
				checkName + "\nFailed saved in file " + name)
		}
	}
}
