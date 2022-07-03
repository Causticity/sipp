// Copyright Raul Vera 2015-2021

package shist

import (
	//"fmt"
	"math/bits"
)
import (
	. "github.com/Causticity/sipp/scomplex"
	. "github.com/Causticity/sipp/simage"
)

// A sparseSippHist is a 2-dimensional histogram of the values in a complex
// gradient image, stored sparsely to conserve memory. This is useful for
// 16-bit and deeper images, as a full flat histogram would require 2^17^2 bins.
type sparseSippHist struct {
	// Embed the core fields
	histCore
	// The histogram data
	sparse map[complex128]uint32
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

func maxSparseHistSize(grad *ComplexImage) int {
	// The maximum size of a sparse histogram is one sparseHistogramEntry per
	// gradient pixel, but Go maps always have a power of 2 number of entries.
	// See the comment for sparseHistogramEntrySize above.
	return int(sparseHistogramEntrySize * upToNextPowerOf2(uint32(len(grad.Pix))))
}

func makeSparseHist(grad *ComplexImage, width, height int) SippHist {
	// A sparse histogram is a map of actually occurring values.
	hist := new(sparseSippHist)
	hist.grad = grad
	hist.width = width
	hist.height = height
	hist.sparse = make(map[complex128]uint32)
	var numUsedBins uint32
	for _, pixel := range grad.Pix {
		v := hist.sparse[pixel]
		if v == 0 {
			// First use of this bin, so count it
			numUsedBins++
		}
		v++
		hist.sparse[pixel] = v
		if v > hist.max {
			hist.max = v
		}
	}
	// numUsedBins is larger, or in the worst case equal, to the number of
	// distinct bin values, so it can be used as the capacity of the slice of
	// distinct bin values.
	hist.bins = make([]BinPair, 0, numUsedBins)
	for _, binval := range hist.sparse {
		addBinsValue(&hist.bins, binval)
	}
	return hist
}

// BinForPixel returns the bin index in the slice returned by Bins for the
// given gradient-image pixel.
func (hist *sparseSippHist) BinForPixel(x, y int) (int) {
	// get the complex value in the image at x, y
	stride := hist.grad.Rect.Dx()
	index := y*stride+x
	pixel := hist.grad.Pix[index]
	// get the value from the map
	val:= hist.sparse[pixel]
	// find the value in the bins slice
	for i, binVal := range hist.bins {
		if val == binVal.BinVal {
			return i
		}
	}
	panic("No bin found for pixel!");
}

// Render renders the histogram by clipping all values to 255. Returns an 8-bit
// grayscale image.
func (hist *sparseSippHist) Render(clip bool) SippImage {
	return nil
}

// RenderSuppressed renders a suppressed version of the histogram and returns
// the result as an 8-bit grayscale image.
func (hist *sparseSippHist) RenderSuppressed() SippImage {
	return nil
}

// RenderSubstitute renders an 8-bit image of the histogram, substituting
// the given value as the pixel value for each corresponding bin value. The
// input slice must be the same length as the slice of bin values returned
// by Bins, and contain new values corresponding to that order.
// This is used to render the delentropy values of the histogram.
func (hist *sparseSippHist) RenderSubstitute(subs []uint8, zeroVal uint8) SippImage {
	return nil
}
