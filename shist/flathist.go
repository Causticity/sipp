// Copyright Raul Vera 2015-2021

package shist

import (
	//"fmt"
	//"image"
	"math"
)

import (
	. "github.com/Causticity/sipp/scomplex"
	. "github.com/Causticity/sipp/simage"
)

// flatSippHist is a 2-dimensional histogram of the values in a complex gradient
// image.
type flatSippHist struct {
	// Embed the core fields.
	histCore
	// The histogram data.
	bin []uint32
	// The index of the histogram bin for each gradient image pixel.
	binIndex []int
}

// flatBinSize is the number of uint32s per histogram entry.
const flatBinSize = 1

// flatSize returns the size of the histogram in uint32s.
func flatSize(width, height int) int {
	// The size of a flat histogram is one flatBinSize per bin.
	return width * height * flatBinSize
}

// Make a flat histogram from the given complex gradient image. The width and
// height (of the histogram, not the image) are passed in to avoid recomputing
// them, as they were needed to decide whether to use this histogram or the
// sparse one.
func makeFlatHist(grad *ComplexImage, width, height int) SippHist {
	hist := new(flatSippHist)
	hist.grad = grad
	hist.width = width
	hist.height = height
	histDataSize := hist.width * hist.height
	hist.bin = make([]uint32, histDataSize)
	hist.binIndex = make([]int, uint32(len(grad.Pix)))
	//fmt.Println("Grad image pixels, and binIndex length:", len(grad.Pix))
	var numUsedBins uint32

	// Walk through the image, computing the bin address from the gradient
	// values, storing the bin address in binIndex, and incrementing the
	// bin. Save the maximum bin value and count the number of actually used
	// bins.
	xoff := int((width-1)/2)
	yoff := int((height-1)/2)
	for i, pixel := range grad.Pix {
		u := int(math.Floor(real(pixel))) + xoff
		v := int(math.Floor(imag(pixel))) + yoff
		hist.binIndex[i] = v*int(hist.width) + u
		if hist.bin[hist.binIndex[i]] == 0 {
			// First use of this bin, so count it
			numUsedBins++
		}
		hist.bin[hist.binIndex[i]]++
		if hist.bin[hist.binIndex[i]] > hist.max {
			hist.max = hist.bin[hist.binIndex[i]]
		}
	}

	// numUsedBins is larger, or in the worst case equal, to the number of
	// distinct bin values, so it can be used as the capacity of the slice of
	// distinct bin values.
	hist.bins = make([]BinPair, 0, numUsedBins)
	for _, binval := range hist.bin {
		hist.bins = addBinsValue(hist.bins, binval)
	}
	//fmt.Println("Histogram complete. Maximum bin value:", hist.Max)
	return hist
}

// BinForPixel returns the bin index in the slice returned by Bins for the
// given gradient-image pixel.
func (hist *flatSippHist) BinForPixel(x, y int) (int) {
	stride := hist.grad.Rect.Dx()
	index := y*stride+x
	//fmt.Printf("index into binIndex for pixel %d, %d is %d\n", x, y, index)
	//fmt.Println("binIndex value at that index is ", hist.binIndex[index])
	val := hist.bin[hist.binIndex[index]]
	//fmt.Println("histogram value is ", val)
	for i, binVal := range hist.bins {
		if val == binVal.BinVal {
			return i
		}
	}
	panic("No bin found for pixel!");
}

// Implement the rowSource interface for rendering
// rowVals returns a slice containing the bin values for one complete row of
// the histogram.
func (hist *flatSippHist) rowVals(y int) []uint32 {
	i := y * hist.width
	return hist.bin[i:i+hist.width]
}

func (hist *flatSippHist) Render(clip bool) SippImage {
	return hist.renderCore(hist, clip)
}

func (hist *flatSippHist) RenderSuppressed() SippImage {
	return hist.renderSuppressedCore(hist)
}

func (hist *flatSippHist) RenderSubstitute(subs []uint8, zeroVal uint8) SippImage {
	return hist.renderSubstituteCore(hist, subs, zeroVal)
}
