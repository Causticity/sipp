// Copyright Raul Vera

// Package shist provides functions for computing a histogram of values of an
// image, and for computing and rendering a 2-dimensional histogram of values of
// a complex or ComplexInt32 gradient image.
package shist

import (
	"fmt"
	"math"
)

import (
	. "github.com/Causticity/sipp/scomplex"
	. "github.com/Causticity/sipp/simage"
)

// A BinPair hold the value of a histogram bin, and the number of times that
// value occurs in the histogram. That makes it a kind of histogram of a
// histogram.
type BinPair struct {
	BinVal 	uint32
	Num		uint32
}

// A SippHist is a 2D histogram.
type SippHist interface {
	// Grad returns the gradient image that this histogram is computed from.
	Grad() (*ComplexImage)
	// Size returns the width and height of the full histogram. ints are used
	// instead of uint32s for compatibility with native Go pixel indexing.
	Size() (int, int)
	// Bins returns a compact slice of BinPairs for the histogram, without
	// duplicates. There is no order specified, but each call to Bins returns
	// the values in the same order.
	Bins() ([]BinPair)
	// BinForPixel returns the index in the slice returned by Bins for the
	// given gradient-image pixel.
	BinForPixel(x, y int) (int)
	// Max returns the maximum bin value that occurs in this histogram
	Max() (uint32)
	// Render returns a rendering of this histogram as an 8-bit image.  If clip
	// is true, values are clipped to 255. If clip is false, values are scaled
	// to 255. TODO: Image dimensions are scaled down for very large histograms.
	Render(clip bool) SippImage
	// RenderSuppressed renders a suppressed version of the histogram and returns
	// the result as an 8-bit grayscale image.
	RenderSuppressed() SippImage
	// RenderSubstitute renders an 8-bit image of the histogram, substituting
	// the given value as the pixel value for each corresponding bin value. The
	// input slice must be the same length as the slice of bin values returned
	// by Bins, and contain new values corresponding to that order.
	// This is used to render the delentropy values of the histogram.
	RenderSubstitute(subs []uint8) SippImage
}

// Internally, a 2D histogram can be either "flat", meaning that storage exists
// for every possible bin, or sparse, meaning that storage is allocated in a map
// based on actually occurring values. Which to use is based on criteria
// described below.
// The above interface hides the distinction from clients.

// The core elements that either version includes.
type histCore struct {
	// A reference to the gradient image we are computing from
	grad *ComplexImage
	// Width and height of the histogram, not the image.
	// These should be odd so that there is always a centre point.
	width, height int
	// The maximum bin value in the histogram.
	max uint32
	// The set of bin values that actually occur, and the number of their
	// occurrences.
	bins []BinPair
	// The inverse of bins. For a given bin value, stores the index in bins.
	// Lazily initialised as needed.
	invertedBins map[uint32]int
}

func (hist *histCore) Grad() (*ComplexImage) {
	return hist.grad
}

func (hist *histCore) Size() (int, int) {
	return hist.width, hist.height
}

func (hist *histCore) Max() (uint32) {
	return hist.max
}

// Compute the width and height of the histogram and the maximum excursion.
// The width and height are twice the maximum excursion on the real and
// imaginary axes, respectively, plus one to ensure that the width and height
// are odd so that there is always a single central bin in both dimensions. The
// returned overall maximum excursion is the maximum of the maximum excursions
// for each axis.
func computeHistSize(grad *ComplexImage) (maxExcursion, width, height int) {
	maxRealExcursion := int(math.Max(math.Abs(grad.MaxRe), math.Abs(grad.MinRe)))
	maxImagExcursion := int(math.Max(math.Abs(grad.MaxIm), math.Abs(grad.MinIm)))
	maxExcursion = int(math.Max(float64(maxRealExcursion), float64(maxImagExcursion)))

	width = maxRealExcursion * 2 + 1 // Ensure both are odd
	height = maxImagExcursion * 2 + 1
	return
}

// All images with maximum excursions less than or equal to minSparseExcursion
// will use the flat version of the histogram, in order to avoid the
// computational overhead of sparse histograms for 8-bit and low-excursion
// 16-bit images, even though the sparse histogram would usually be smaller.
// Note that excursion, the distance from 0 along either axis of the complex
// plane, is always positive.
// TODO: Later, this will be used to scale down sparse histograms for rendering,
// so that histogram renderings will never be more that double this value on a
// side.
const minSparseExcursion = 1024

// Hist computes the 2D histogram from the given gradient image.
func Hist(grad *ComplexImage) (hist SippHist) {
	maxExcursion, width, height := computeHistSize(grad)
	// The following sizes are number of uint32s for the histogram.
	flatHistSize := flatSize(width, height)
	maxSparseSize := maxSparseHistSize(grad)

	//fmt.Println("flat histogram width, height: ", width, height)
	//fmt.Println("maxSparseSize:", maxSparseSize, ", flatHistSize:", flatHistSize)
	if maxExcursion > minSparseExcursion && maxSparseSize < flatHistSize {
		// Use a sparse histogram
		fmt.Println("Using sparse histogram")
		hist = makeSparseHist(grad, width, height)
	} else {
		// Use a flat histogram
		fmt.Println("Using flat histogram")
		hist = makeFlatHist(grad, width, height)
	}
	return
}

// supScale determines a scale factor that is the ratio of the distance to
// the given x, y from the given centre, over the given maximum distance.
// This is used for rendering suppressed images of the histogram.
func supScale(x, y, centx, centy int, maxDist float64) float64 {
	xdist := float64(x - centx)
	ydist := float64(y - centy)
	hyp := math.Hypot(xdist, ydist)
	return (hyp / maxDist)
}
