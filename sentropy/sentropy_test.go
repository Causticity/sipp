// Copyright Raul Vera 2020

// Tests for package sentropy.

package sentropy

import (
	"image"
	_ "image/png"
	"path/filepath"
	"reflect"
	"testing"
)

import (
	. "github.com/Causticity/sipp/scomplex"
	. "github.com/Causticity/sipp/shist"
	. "github.com/Causticity/sipp/simage"
	. "github.com/Causticity/sipp/sipptesting"
	. "github.com/Causticity/sipp/sipptesting/sipptestcore"
)

var smallPicEntropy = 4.0

var smallPicEntropyImage = []uint8{
	255, 255, 255, 255,
	255, 255, 255, 255,
	255, 255, 255, 255,
	255, 255, 255, 255,
}

var smallPic16Entropy = 4.0

var cosxCosyTinyEntropy = 5.211838049805185

var cosxCosyTinyEntropyImage = []uint8{
	154, 255, 255, 154, 154, 255, 255, 255, 255, 154, 154, 255, 255, 255, 255, 154, 154, 255, 255, 154,
	255, 154, 208, 154, 255, 255, 124, 208, 154, 255, 255, 154, 208, 182, 255, 255, 154, 208, 154, 255,
	255, 208, 124, 255, 154, 154, 208, 154, 208, 255, 255, 208, 154, 208, 154, 154, 255, 52, 208, 255,
	154, 154, 255, 124, 255, 255, 154, 208, 124, 255, 255, 124, 208, 154, 255, 255, 124, 255, 154, 154,
	154, 255, 154, 255, 154, 154, 255, 154, 255, 255, 255, 255, 154, 255, 154, 154, 255, 154, 255, 154,
	255, 255, 154, 255, 154, 154, 255, 154, 255, 154, 154, 255, 154, 255, 154, 154, 255, 154, 255, 255,
	255, 182, 91, 154, 255, 255, 124, 255, 154, 154, 154, 154, 255, 52, 255, 255, 154, 208, 182, 255,
	255, 91, 154, 91, 154, 154, 255, 124, 208, 255, 255, 91, 124, 255, 154, 154, 208, 154, 208, 255,
	255, 154, 91, 182, 255, 255, 154, 208, 154, 255, 255, 154, 91, 154, 255, 255, 182, 208, 154, 255,
	154, 255, 255, 255, 255, 154, 154, 255, 255, 154, 154, 255, 255, 154, 154, 255, 255, 255, 255, 154,
	154, 255, 255, 255, 255, 154, 154, 255, 255, 154, 154, 255, 255, 154, 154, 255, 255, 255, 255, 154,
	255, 154, 208, 182, 255, 255, 154, 208, 154, 255, 255, 154, 208, 154, 255, 255, 124, 208, 154, 255,
	255, 208, 154, 208, 154, 154, 255, 52, 208, 255, 255, 208, 124, 255, 154, 154, 208, 154, 208, 255,
	255, 124, 208, 154, 255, 255, 124, 255, 154, 154, 154, 154, 255, 124, 255, 255, 154, 208, 124, 255,
	255, 255, 154, 255, 154, 154, 255, 154, 255, 154, 154, 255, 154, 255, 154, 154, 255, 154, 255, 255,
	154, 255, 154, 255, 154, 154, 255, 154, 255, 255, 255, 255, 154, 255, 154, 154, 255, 154, 255, 154,
	154, 154, 255, 52, 255, 255, 154, 208, 182, 255, 255, 182, 91, 154, 255, 255, 124, 255, 154, 154,
	255, 91, 124, 255, 154, 154, 208, 154, 208, 255, 255, 91, 154, 91, 154, 154, 255, 124, 208, 255,
	255, 154, 91, 154, 255, 255, 182, 208, 154, 255, 255, 154, 91, 182, 255, 255, 154, 208, 154, 255,
	154, 255, 255, 154, 154, 255, 255, 255, 255, 154, 154, 255, 255, 255, 255, 154, 154, 255, 255, 154,
}

var expectedDelentropy float64 = 6.775012499324653
var expectedMaxDelentropy float64 = 0.12179180114985422
var expectedDelentropyArray = []float64{
	0.023534224451211,
	0.04152828269743585,
	0.09824198104431049,
	0.12179180114985422,
	0.0719762329848994,
	0.0855114533517979,
}

var cosxCosyTinyDelentropyImageArray = []uint8{
	205, 150, 86, 86, 86, 86, 86, 150, 150, 49, 86, 205, 205, 86, 205, 86, 86, 205, 205,
	205, 205, 150, 150, 86, 205, 86, 255, 150, 86, 150, 205, 86, 205, 86, 86, 150, 86, 205,
	86, 86, 86, 86, 86, 86, 255, 86, 86, 49, 86, 86, 205, 86, 205, 205, 205, 150, 86,
	86, 205, 205, 205, 179, 150, 86, 205, 86, 86, 86, 150, 86, 205, 179, 205, 205, 86, 86,
	205, 86, 205, 179, 179, 179, 86, 86, 86, 49, 86, 86, 86, 179, 179, 179, 205, 86, 205,
	86, 150, 205, 205, 179, 205, 86, 150, 86, 86, 86, 205, 150, 150, 179, 205, 86, 205, 86,
	150, 150, 205, 205, 205, 205, 86, 86, 86, 49, 205, 86, 255, 150, 86, 86, 205, 86, 205,
	86, 205, 150, 150, 86, 205, 86, 205, 150, 86, 86, 255, 86, 205, 86, 86, 150, 205, 205,
	86, 86, 150, 86, 205, 86, 86, 205, 205, 179, 150, 86, 205, 86, 86, 86, 150, 150, 86,
	49, 205, 86, 205, 49, 205, 86, 205, 179, 179, 179, 86, 49, 86, 49, 86, 49, 86, 49,
	86, 205, 205, 86, 205, 86, 86, 205, 205, 179, 205, 150, 86, 86, 86, 86, 86, 150, 150,
	150, 205, 86, 205, 86, 86, 150, 86, 205, 205, 205, 205, 150, 150, 86, 205, 86, 255, 150,
	86, 86, 205, 86, 205, 205, 205, 150, 86, 86, 86, 86, 86, 86, 86, 86, 255, 86, 86,
	86, 150, 86, 205, 179, 205, 205, 86, 86, 205, 86, 205, 205, 205, 179, 150, 86, 205, 86,
	86, 86, 86, 179, 179, 179, 205, 86, 205, 49, 205, 86, 205, 179, 179, 179, 86, 86, 86,
	86, 205, 150, 150, 179, 205, 86, 205, 86, 205, 86, 150, 205, 205, 179, 205, 86, 150, 86,
	205, 86, 255, 150, 86, 86, 205, 86, 205, 86, 150, 150, 205, 205, 205, 205, 86, 86, 86,
	86, 255, 86, 205, 86, 86, 150, 205, 205, 205, 86, 205, 150, 150, 86, 205, 86, 205, 150,
	150, 86, 205, 86, 86, 86, 150, 150, 86, 49, 86, 86, 150, 86, 205, 86, 86, 205, 205,
}

var CosxCosyTinyDelentropyImage = image.Gray{
	cosxCosyTinyDelentropyImageArray,
	19,
	image.Rectangle{image.Point{0, 0}, image.Point{19, 19}},
}

var sgrayCosxCosyTinyDelentropy *SippGray

func init() {
	sgrayCosxCosyTinyDelentropy = new(SippGray)
	sgrayCosxCosyTinyDelentropy.Gray = &CosxCosyTinyDelentropyImage
}

func TestConventionalEntropy(t *testing.T) {
	ent := Entropy(SgrayCosxCosyTiny)
	if ent.Entropy != cosxCosyTinyEntropy {
		t.Errorf("Error: Incorrect conventional entropy: expected %v, got %v",
			cosxCosyTinyEntropy, ent.Entropy)
	}
	entIm := ent.EntropyImage()
	if !reflect.DeepEqual(entIm.Pix(), cosxCosyTinyEntropyImage) {
		t.Error("Error: Entropy image incorrect. Expected:" +
			GrayArrayToString(cosxCosyTinyEntropyImage, CosxCosyTinyStride) + "Got:" +
			GrayArrayToString(entIm.Pix(), CosxCosyTinyStride))
	}
	ent16 := Entropy(Sgray16)
	if ent16.Entropy != smallPic16Entropy {
		t.Errorf("Error: Incorrect conventional entropy: expected %v, got %v",
			smallPic16Entropy, ent16.Entropy)
	}
	entIm16 := ent16.EntropyImage()
	if !reflect.DeepEqual(entIm16.Pix(), smallPicEntropyImage) {
		t.Error("Error: Entropy image incorrect. Expected:" +
			GrayArrayToString(smallPicEntropyImage, 4) + "Got:" +
			GrayArrayToString(entIm16.Pix(), 4))
	}
}

type entropyTest struct {
	name                    string
	grad                    []complex128
	stride                  int
	maxDelentropy           float64
	delentropyArray         []float64
	delentropy              float64
	histDelentropyImageName string
	delentropyImage         *SippGray
}

func TestDelentropy(t *testing.T) {
	var tests = []entropyTest{
		{
			"CosxCosyTinyGrad",
			CosxCosyTinyGrad,
			CosxCosyTinyStride,
			expectedMaxDelentropy,
			expectedDelentropyArray,
			expectedDelentropy,
			"cosxcosy_tiny_hist_delent.png",
			sgrayCosxCosyTinyDelentropy,
		},
	}
	for _, test := range tests {
		hist := Hist(FromComplexArray(test.grad, test.stride-1))
		dent := Delentropy(hist)
		if dent.hist != hist {
			t.Errorf("Error: SippDelentropy for %s has incorrect hist, expected %v, got %v",
				test.name, hist, dent.hist)
		}
		if !reflect.DeepEqual(test.delentropyArray, dent.binDelentropy) {
			t.Errorf("Error: delentropy array for %s incorrect. Expected %v, got %v",
				test.name, test.delentropyArray, dent.binDelentropy)
		}
		if test.maxDelentropy != dent.maxBinDelentropy {
			t.Errorf("Error: maxBinDelentropy for %s incorrect. Expected %v, got %v",
				test.name, test.maxDelentropy, dent.maxBinDelentropy)
		}
		if dent.Delentropy != test.delentropy {
			t.Errorf("Error: delentropy for %s incorrect. Expected %v, got %v",
				test.name, test.delentropy, dent.Delentropy)
		}

		histDentImage := dent.HistDelentropyImage()
		checkName := filepath.Join(TestDir, test.histDelentropyImageName)
		check, err := Read(checkName)
		if err != nil {
			t.Errorf("Error reading histogram delentropy check image: %v\n", test.histDelentropyImageName)
		}
		if !reflect.DeepEqual(histDentImage.Pix(), check.Pix()) {
			// Write out the check image and report names of mismatched files
			name := SaveFailedSimage(checkName, histDentImage)
			t.Errorf("Error: histogram delentropy image does not match expected.\nExpected in file " +
				checkName + "\nFailed saved in file " + name)
		}

		delentImage := dent.DelEntropyImage()
		if !reflect.DeepEqual(delentImage.Pix(), test.delentropyImage.Pix()) {
			t.Errorf("Error: gradient delentropy image incorrect. Expected %v, got %v\n",
				test.delentropyImage.Pix(), delentImage.Pix())
		}
	}
}
