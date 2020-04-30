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
	. "github.com/Causticity/sipp/sgrad"
	. "github.com/Causticity/sipp/simage"
    . "github.com/Causticity/sipp/sipptesting/sipptestcore"
	. "github.com/Causticity/sipp/sipptesting"
)

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

func TestConventionalEntropy(t *testing.T) {
    ent, entIm := Entropy(SgrayCosxCosyTiny)
     if ent != CosxCosyTinyEntropy {
        t.Errorf("Error: Incorrect conventional entropy: expected %v, got %v",
                 CosxCosyTinyEntropy, ent)
    } 
    if !reflect.DeepEqual(entIm.Pix(), CosxCosyTinyEntropyImage) {
        t.Error("Error: Entropy image incorrect. Expected:" +
            GrayArrayToString(CosxCosyTinyEntropyImage, CosxCosyTinyStride) + "Got:" +
            GrayArrayToString(entIm.Pix(), CosxCosyTinyStride))
    }
    ent16, entIm16 := Entropy(Sgray16)
    if ent16 != SmallPic16Entropy {
        t.Errorf("Error: Incorrect conventional entropy: expected %v, got %v",
                 SmallPic16Entropy, ent16)
    }
    if !reflect.DeepEqual(entIm16.Pix(), SmallPicEntropyImage) {
        t.Error("Error: Entropy image incorrect. Expected:" +
            GrayArrayToString(SmallPicEntropyImage, 4) + "Got:" +
            GrayArrayToString(entIm16.Pix(), 4))
    }
}

type histTest struct {
    name string
    grad []complex128
    k uint16
    stride int
    maxMod float64
    count int
    binIndex []int
    binVals []uint32
    maxBinVal uint32
    delentropy float64
    maxdelentropy float64
    delentropyArray [] float64
    histDelentropyImageName string
    delentropyImage *SippGray
    suppressedImageName string
    renderedHistogramName string
}

func TestHist(t *testing.T) {
    var tests = []histTest {
        {
            "CosxCosyTinyGrad K = 0",
            CosxCosyTinyGrad,
            0,
            CosxCosyTinyStride,
            CosxCosyTinyGradMaxMod,
            ExpectedCosxCosyTinyNonZeroHistCount,
            CosxCosyTinyBinIndex,
            CosxCosyTinyBinVals,
            ExpectedMax,
            ExpectedDelentropy,
            ExpectedMaxDelentropy,
            ExpectedDelentropyArray,
            "cosxcosy_tiny_hist_delent.png",
            SgrayCosxCosyTinyDelentropy,
            "cosxcosy_tiny_hist_sup.png",
            "cosxcosy_tiny_hist.png",
        },
        {
            "CosxCosyTinyGrad K = 255",
            CosxCosyTinyGrad,
            255,
            CosxCosyTinyStride,
            CosxCosyTinyGradMaxMod,
            ExpectedCosxCosyTinyNonZeroHistCount,
            CosxCosyTinyBinIndexk255,
            CosxCosyTinyBinVals,
            ExpectedMax,
            ExpectedDelentropy,
            ExpectedMaxDelentropy,
            ExpectedDelentropyArray,
            "cosxcosy_tiny_k255_hist_delent.png",
            SgrayCosxCosyTinyDelentropy,
            "cosxcosy_tiny_k255_hist_sup.png",
            "cosxcosy_tiny_k255_hist.png",
        },
        // TODO 16-bit tests. Depend on grad 16-bit tests
        // {
            // "CosxCosyTinyGrad16",
            // CosxCosyTinyGrad16,
            // 0,
            // CosxCosyTiny16Stride,
            // CosxCosyTinyGrad16MaxMod,
            // ExpectedCosxCosyTiny16NonZeroHistCount,
            // CosxCosyTiny16BinIndex,
            // CosxCosyTiny16BinVals,
            // ExpectedMax16,
            // ExpectedDelentropy16,
            // ExpectedMaxDelentropy16,
            // ExpectedDelentropyArray16,
            // "cosxcosy_tiny_16_hist_delent.png",
            // SgrayCosxCosyTinyDelentropy16,
            // "cosxcosy_tiny_16_hist_sup.png",
            // "cosxcosy_tiny_16_hist.png",
        // },
    }
    for _, test := range tests {
        grad := FromComplexArray(test.grad, test.stride-1)
        hist := Hist(grad, test.k)
        if hist.grad != grad {
            t.Errorf("Error: SippHist for %s has incorrect grad, expected %v, got %v",
                test.name, grad, hist.grad)
        }
        xpctK := test.k
        if test.k == 0 {
            xpctK = uint16(test.maxMod) + kMargin
        }
        if hist.k != xpctK {
            t.Errorf("Error: K for %s histogram incorrect. Expected %v, got %v", 
                test.name, xpctK, hist.k)
        }
        count := 0
        for _, val := range hist.bin {
            if val != 0 {
                count++
            }
        }
        xpctCnt := test.count
        if count != xpctCnt {
            t.Errorf("Error: Histogram for %s has incorrect number of non-zero entries: expected %v, got %v", 
                test.name, xpctCnt, count) 
        }
        if !reflect.DeepEqual(hist.binIndex, test.binIndex) {
            t.Errorf("Error: hist.binIndex for %s incorrect, expected\n%v\n got\n%v\n", 
                test.name, test.binIndex, hist.binIndex)
        }
        for i, val := range hist.binIndex {
            if hist.bin[val] != test.binVals[i] {
                t.Errorf("Error: histogram value for %s incorrect, expected %v, got %v",
                    test.name, test.binVals[i], hist.bin[val])
            }
        }
        if hist.max != test.maxBinVal {
            t.Errorf("Error: hist.max for %s incorrect. Expected %v, got %v",
                test.name, test.maxBinVal, hist.max)
        }
        histDentImage := hist.HistDelentropyImage()
        if histDentImage != nil {
            t.Error("HistDelentropyImage returned non-nil before Delentropy called")
        }
        dent := hist.Delentropy()
        if dent != test.delentropy {
            t.Errorf("Error: delentropy for %s incorrect. Expected %v, got %v",
                test.name, test.delentropy, dent)
        }
        if test.maxdelentropy != hist.maxDelentropy {
            t.Errorf("Error: maxdelentropy for %s incorrect. Expected %v, got %v",
                test.name, test.maxdelentropy, hist.maxDelentropy)
        }
        if !reflect.DeepEqual(test.delentropyArray, hist.delentropy) {
            t.Errorf("Error: delentropy array for %s incorrect. Expected %v, got %v",
                test.name, test.delentropyArray, hist.delentropy)
        }
        histDentImage = hist.HistDelentropyImage()
        check, err := Read(filepath.Join(TestDir, test.histDelentropyImageName))
        if err != nil {
            t.Errorf("Error reading histogram delentropy check image: %v\n", test.histDelentropyImageName)
        }
        if !reflect.DeepEqual(histDentImage.Pix(), check.Pix()) {
            t.Errorf("Error: histogram delentropy image incorrect. Expected %v, got%v\n",
                check.Pix(), histDentImage.Pix())
        }
        delentImage := hist.DelEntropyImage()
        if !reflect.DeepEqual(delentImage.Pix(), test.delentropyImage.Pix()) {
            t.Errorf("Error: gradient delentropy image incorrect. Expected %v, got %v\n",
                test.delentropyImage.Pix(), delentImage.Pix())
        }
        supp := hist.RenderSuppressed()
        check, err = Read(filepath.Join(TestDir, test.suppressedImageName))
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