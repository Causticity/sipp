// Copyright Raul Vera 2015-2020

// Package shist provides functions for computing a histogram of values of an
// image, and for computing and rendering a 2-dimensional histogram of values of
// a complex or ComplexInt32 gradient image.
package shist

import (
	"fmt"
	"image"
	"math"
)

import (
	. "github.com/Causticity/sipp/scomplex"
	. "github.com/Causticity/sipp/simage"
)

// SippHist is a 2-dimensional histogram of the values in a complex gradient
// image.
type SippHist struct {
	// A reference to the gradient image we are computing from
	Grad *ComplexImage
	// These should be odd so that there is always a centre point.
	width, height uint32
	// The histogram data.
	Bin []uint32
	// The index of the histogram bin for each gradient image pixel.
	BinIndex []int
	// The maximum bin value in the histogram.
	Max uint32
	// A suppressed version of the histogram, stored as floats for subsequent
	// computation.
	suppressed []float64
	// The maximum suppressed value, stored as a float for subsequent
	// computation.
	maxSuppressed float64
}

func (hist *SippHist) Size() (uint32, uint32) {
	return hist.width, hist.height
}

const greyHistSize8BPP = 256
const greyHistSize16BPP = 65536

// GreyHist computes a 1D histogram of the greyscale values in the image.
func GreyHist(im SippImage) (hist []uint32) {
	histSize := greyHistSize8BPP
	is16 := false
	if im.Bpp() == 16 {
		histSize = greyHistSize16BPP
		is16 = true
	}

	//fmt.Println("GreyHist histogram size is ", histSize)

	hist = make([]uint32, histSize)
	imPix := im.Pix()
	for y := 0; y < im.Bounds().Dy(); y++ {
		for x := 0; x < im.Bounds().Dx(); x++ {
			index := im.PixOffset(x, y)
			var val uint16 = uint16(imPix[index])
			if is16 {
				val = val<<8 | uint16(imPix[index+1])
			}
			hist[val]++
		}
	}
	return
}

type sparseHistogramEntry struct {
	// The coordinates in histogram space for this bin
	x, y uint32
	// The bin value, i.e. the number of pixels with this gradient value
	binVal uint32
}

// number of uint32s per bin for each of the flat or sparse versions
const sparseHistogramEntrySize = 3
const flatBinSize = 1

// All images with excursions <= this will use the flat version, in order to
// avoid the computational overhead of sparse histograms for 8-bit and
// low-excursion 16-bit images, even though the sparse histogram would usually
// be smaller.
const minSparseExcursion = 256

// Hist computes the 2D histogram from the given gradient image.
func Hist(grad *ComplexImage) (hist *SippHist) {
	hist = new(SippHist)
	hist.Grad = grad

	// The size of a flat histogram is one flatBinSize per bin. The number
	// of bins is the product of the histogram width and height. The width and
	// height are twice the maximum excursion on the real and imaginary axes,
	// respectively, plus one to ensure that the width and height are odd so
	// that there is always a single central bin in both dimensions.
	maxRealExcursion := uint32(math.Max(math.Abs(grad.MaxRe), math.Abs(grad.MinRe)))
	maxImagExcursion := uint32(math.Max(math.Abs(grad.MaxIm), math.Abs(grad.MinIm)))
	maxExcursion := uint32(math.Max(float64(maxRealExcursion), float64(maxImagExcursion)))

	width := maxRealExcursion * 2 + 1 // Ensure both are odd
	height := maxImagExcursion * 2 + 1
	flatHistSize := width * height * flatBinSize

	nPix := len(grad.Pix)

	// Compute the size of the regular histogram and the maximum size of a sparse
	// histogram and use the smaller version

	// The maximum size of a sparse histogram is one sparseHistogramEntry per
	// gradient pixel. The size is the number of uint32s.
	maxSparseSize := uint32(sparseHistogramEntrySize * nPix)

	//fmt.Println("flat histogram width, height: ", width, height)
	//fmt.Println("maxSparseSize:", maxSparseSize, ", flatHistSize:", flatHistSize)
	if maxExcursion > minSparseExcursion && maxSparseSize < flatHistSize {
		// Use a sparse histogram
		fmt.Println("Would use sparse histogram, but unimplemented")
		panic(1)
	} else {
		// Use a flat histogram
		//fmt.Println("Using flat histogram")

		hist.width = width
		hist.height = height

		histDataSize := int(hist.width) * int(hist.height) // Always odd
		hist.Bin = make([]uint32, histDataSize)
		hist.BinIndex = make([]int, nPix)

		// Walk through the image, computing the bin address from the gradient
		// values storing the bin address in BinIndex and incrementing the bin.
		// Save the maximum bin value as well.
		for i, pixel := range grad.Pix {
			u := int(math.Floor(real(pixel))) + int(maxRealExcursion)
			v := int(math.Floor(imag(pixel))) + int(maxImagExcursion)
			hist.BinIndex[i] = v*int(hist.width) + u
			hist.Bin[hist.BinIndex[i]]++
			if hist.Bin[hist.BinIndex[i]] > hist.Max {
				hist.Max = hist.Bin[hist.BinIndex[i]]
			}
		}
		//fmt.Println("Histogram complete. Maximum bin value:", hist.Max)
	}
	return
}

// supScale determines a scale factor that is the ratio of the distance to
// the given x, y from the centre, over the given maximum distance.	We assume
// that width and height are odd, so that there is an exact centre.
func supScale(x, y, width, height int, maxDist float64) float64 {
	xdist := float64(x - (width-1)/2)
	ydist := float64(y - (height-1)/2)
	hyp := math.Hypot(xdist, ydist)
	return (hyp / maxDist)
}

// Suppress suppresses the spike near the origin of the histogram by scaling
// the values in the histogram by a facter determined by their distance from the
// origin.
func (hist *SippHist) suppress() {
	if hist.suppressed != nil {
		return
	}
	size := int(hist.width) * int(hist.height)
	hist.suppressed = make([]float64, size)
	var index uint32 = 0
	hist.maxSuppressed = 0
	for y := 0; y < int(hist.height); y++ {
		for x := 0; x < int(hist.width); x++ {
			sscale := supScale(x, y, int(hist.width), int(hist.height), hist.Grad.MaxMod)
			hist.suppressed[index] = float64(hist.Bin[index]) * sscale
			if hist.suppressed[index] > hist.maxSuppressed {
				hist.maxSuppressed = hist.suppressed[index]
			}
			index++
		}
	}
	//fmt.Println("Distance suppression complete; max suppressed value:", hist.maxSuppressed)
}

// RenderSuppressed renders a suppressed version of the histogram and returns
// the result as an 8-bit grayscale image.
func (hist *SippHist) RenderSuppressed() SippImage {
	// Here we will generate an 8-bit output image of the same size as the
	// histogram, scaled to use the full dynamic range of the image format.
	hist.suppress()
	width, height := hist.Size()
	var scale float64 = 255.0 / hist.maxSuppressed
	//fmt.Println("Suppressed Render scale factor:", scale)
	rnd := new(SippGray)
	rnd.Gray = image.NewGray(image.Rect(0, 0, int(width), int(height)))
	rndPix := rnd.Pix()
	for index, val := range hist.suppressed {
		rndPix[index] = uint8(val * scale)
	}
	return rnd
}

// Render renders the histogram by clipping all values to 255. Returns an 8-bit
// grayscale image.
func (hist *SippHist) Render() SippImage {
	// Here we will generate an 8-bit output image of the same size as the
	// histogram, clipped to 255.
	width, height := hist.Size()
	//var scale float64 = 255.0 / float64(hist.Max)
	//fmt.Println("Render scale factor:", scale)
	rnd := new(SippGray)
	rnd.Gray = image.NewGray(image.Rect(0, 0, int(width), int(height)))
	rndPix := rnd.Pix()
	for index, val := range hist.Bin {
		if val > 255 {
			val = 255
		}
		rndPix[index] = uint8(val)
	}
	return rnd
}
