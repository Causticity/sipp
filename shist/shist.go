// Copyright Raul Vera 2015-2020

// Package shist provides functions for computing a histogram of values of an
// image, and for computing and rendering a 2-dimensional histogram of values of
// a complex or ComplexInt32 gradient image.
package shist

import (
	"fmt"
	"image"
	"math"
	"math/bits"
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

// sparseHistogramEntrySize is the number of uint32s per histogram entry.
// The size in bytes (or uint32s) of a Go map is not easy to determine, but the
// number of buckets is always a power of 2, so as a rough estimate, we'll take
// the minimum size of an entry to be the size of the complex128 index (4 uint32s)
// plus the count (1 uint32) plus a 64-bit pointer for overhead (2 uint32s). The
// last of these is just a wild guess. Then we multiply this entry size by the
// number of pixels rounded up to the next power of 2 to get an estimate of the
// sparse histogram size.
const sparseHistogramEntrySize = 4 + 1 + 2 // See above
const flatBinSize = 1

// Return the next power of 2 higher than the input, or panic if the input is 0,
// 1, or would overflow, as none of these should ever occur.
func upToNextPowerOf2(n uint32) uint32 {
	// An input of 0 or 1 is invalid, but should never be sent in, so panic.
	if n == 0 || n == 1 {
		panic(1)
	}
	// Next check if the input already is a power of 2.
	if bits.OnesCount32(n) == 1 {
		return n
	}
	// Now check that the highest bit isn't already set, because if it is then
	// we would overflow. This means we have more than 2 billion pixels, which
	// would be problematic for plenty of other reasons, so here we assume that
	// this is an error in the caller and just panic.
	lz := bits.LeadingZeros32(n)
	if lz == 0 {
		panic(1)
	}
	// Otherwise return a number with a single bit set one position higher than
	// the highest input bit set.
	return 1 << (32 - lz)
}

// All images with excursions <= this will use the flat version, in order to
// avoid the computational overhead of sparse histograms for 8-bit and
// low-excursion 16-bit images, even though the sparse histogram would usually
// be smaller. Note that excursion, the distance from 0 along either axis of the
// complex plane, is always positive.
// Later, this will be used to scale down sparse histograms for rendering, so
// that histogram renderings will never be more that double this value on a side.
const minSparseExcursion = 1024

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

	hist.width = maxRealExcursion * 2 + 1 // Ensure both are odd
	hist.height = maxImagExcursion * 2 + 1
	flatHistSize := hist.width * hist.height * flatBinSize

	nPix := uint32(len(grad.Pix))

	// Compute the size of the regular histogram and the maximum size of a sparse
	// histogram and use the smaller version

	// The maximum size of a sparse histogram is one sparseHistogramEntry per
	// gradient pixel, but Go maps always have a power of 2 number of entries.
	// See the comment for sparseHistogramEntrySize above.
	numMapEntries := upToNextPowerOf2(nPix)
	maxSparseSize := uint32(sparseHistogramEntrySize * numMapEntries)

	//fmt.Println("flat histogram width, height: ", width, height)
	//fmt.Println("maxSparseSize:", maxSparseSize, ", flatHistSize:", flatHistSize)
	if maxExcursion > minSparseExcursion && maxSparseSize < flatHistSize {
		// Use a sparse histogram
		fmt.Println("Using sparse histogram")
		// TODO: No other code uses this yet.
		// A sparse histogram is a map of actually occurring values.
		sparse := make(map[complex128]uint32)
		for _, pixel := range grad.Pix {
			v := sparse[pixel]
			v++ // v is 0 for the empty initial case, so this always works
			sparse[pixel] = v
			if v > hist.Max {
				hist.Max = v
			}
		}
	} else {
		// Use a flat histogram
		fmt.Println("Using flat histogram")

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
