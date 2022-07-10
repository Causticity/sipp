// Copyright Raul Vera

// Package shist provides functions for computing a histogram of values of an
// image, and for computing and rendering a 2-dimensional histogram of values of
// a complex or ComplexInt32 gradient image.
package shist

import (
	"image"
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
	// Max returns the maximum bin value that occurs in this histogram
	Max() (uint32)
	// Bins returns a compact slice of BinPairs for the histogram, without
	// duplicates. There is no order specified, but each call to Bins returns
	// the values in the same order.
	Bins() ([]BinPair)
	// BinForPixel returns the index in the slice returned by Bins for the
	// given gradient-image pixel.
	BinForPixel(x, y int) (int)
	// Render returns a rendering of this histogram as an 8-bit image.  If clip
	// is true, values are clipped to 255. If clip is false, values are scaled
	// to 255.
	Render(clip bool) SippImage
	// RenderSuppressed renders a suppressed version of the histogram and returns
	// the result as an 8-bit grayscale image. "Suppressed" here means that the
	// bin value is scaled by the ratio of the distance of the bin from 0 over
	// the maximum distance. Hence the centre is reduced to 0, the maximum
	// value remains unmodified, and all values in between will be scaled
	// linearly between them.
	RenderSuppressed() SippImage
	// RenderSubstitute renders an 8-bit image of the histogram, substituting
	// the given value as the pixel value for each corresponding bin value. The
	// input slice must be the same length as the slice of bin values returned
	// by Bins, and contain new values corresponding to that order.
	// Note that as the slice returned by Bins() does not include 0 values,
	// the value to be used for empty bins must be supplied.
	// This is used to render point transforms of the histogram, such as the
	// delentropy values.
	RenderSubstitute(subs []uint8, zeroVal uint8) SippImage
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

func (hist *histCore) Bins() ([]BinPair) {
	return hist.bins
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

// addBinsValue adds to the bins slice, either by incrementing the Num of the
// relevant entry if the value is already in the slice, or appending a new entry
// if it isn't. Returns the slice, or a new one if append is used.
func addBinsValue(bins []BinPair, binval uint32) []BinPair{
	var found bool
	if binval != 0 {
		found = false
		//fmt.Printf("looking for binval %d\n", binval)
		for i, pair := range bins {
			if binval == pair.BinVal {
				//fmt.Println("Found it, incrementing pair.num")
				bins[i].Num++
				found = true
			}
		}
		if found == false {
			//fmt.Printf("Appending pair for binval %d\n", binval)
			bins = append(bins, BinPair{binval, 1})
		}
	}
	return bins
}

// The maximum size of either dimension of a rendering of a histogram. If
// either excursion is such that the image would be larger than this in either
// dimension, then the histogram is scaled, equally in both dimensions to
// preserve aspect ratio, so that the larger dimension is this size.
const maxRenderExtent = 4096

// renderInto creates and returns a new 8-bit Sippimage correctly sized for the
// histogram, along with the scale factor used and the pixel scale factor
// mapping the maximum bin value to 255. The correct size is determined by
// taking maxRenderExtent into account. If the histogram is already
// smaller in both dimensions, the existing width and height are used.
// Otherwise, the larger dimension is scaled down to maxRenderExtent and the
// other dimension is scaled by the same factor, preserving the aspect ratio.
func (hist *histCore) renderInto() (rnd *SippGray, scale, pixScale float64) {
	rnd = new(SippGray)
	var width, height int
	width, height = hist.Size()
	if (width <= maxRenderExtent && height <= maxRenderExtent) {
		scale = 1.0
	} else if hist.width > maxRenderExtent {
		width = maxRenderExtent
		scale = float64(hist.width)/float64(maxRenderExtent)
		height = int(float64(hist.height)/scale)
	} else {
		height = maxRenderExtent
		scale = float64(hist.height)/float64(maxRenderExtent)
		width = int(float64(hist.width)/scale)
	}
	rnd.Gray = image.NewGray(image.Rect(0, 0, int(width), int(height)))
	pixScale = 255.0 / float64(hist.max)
	return
}

// A rowSource provides a method for returning a complete row of the histogram.
// Both histogram types, flat and sparse, implement this method, which is used
// for rendering.
type rowSource interface {
	// rowVals returns a slice containing the bin values for one complete row of
	// the histogram.
	rowVals(y int) []uint32
}

// renderCore renders the histogram into an 8-bit grayscale image. The calling
// histogram provides itself as the row source. If clip is true, values are
// clipped to 255. If clip is false, values are scaled to 255.
func (hist *histCore) renderCore(rs rowSource, clip bool) SippImage {
	rnd, scale, pixScale := hist.renderInto()
	rndPix := rnd.Pix()
	//fmt.Println("Render pixel scale factor:", pixScale)
	if scale == 1.0  {
		for row := 0; row < hist.height; row++ {
			histRow := rs.rowVals(row)
			for index, val := range histRow {
				if clip {
					if val > 255 {
						val = 255
					}
				} else {
					val = uint32(math.Round(float64(val) * pixScale))
				}
				rndPix[row*hist.width+index] = uint8(val)
			}
		}
	} else {
		// Render by rows to ensure that the intermediate fits in memory
		// Use a block as wide as the output image but only as tall as one
		// vertical filter. Fill this block, generate the output row, then
		// roll the block.
		// Allocate the intermediate block as a ring buffer of rows
		// precompute both filters
		// fill the first row of the ring buffer by filtering horizontally
		// for each ouput row
			// fill the remaining rows of the ring buffer by filtering horizontally
			// filter the buffer vertically into the output row
			// discard all the rows except the last (if the filter is non-zero!)
			// make the last row the first
	}
	return rnd
}

// supScale determines a scale factor that is the ratio of the distance to
// the given x, y from the given centre, over the given maximum distance.
// This is used for rendering suppressed images of the histogram.
func supScale(x, y int, centx, centy, maxDist float64) float64 {
	xdist := float64(x) - centx
	ydist := float64(y) - centy
	hyp := math.Hypot(xdist, ydist)
	return (hyp / maxDist)
}

// RenderSuppressed renders a suppressed version of the histogram and returns
// the result as an 8-bit grayscale image.
func (hist *histCore) renderSuppressedCore(rs rowSource) SippImage {
	rnd, filterScale, _ := hist.renderInto()
	rndWidth := rnd.Bounds().Dx()
	rndHeight := rnd.Bounds().Dy()
	var suppressed []float64
	var maxSuppressed float64
	size := int(rndWidth) * int(rndHeight)
	suppressed = make([]float64, size)
	centx := (float64(hist.width)-1)/2
	centy := (float64(hist.height)-1)/2

	if filterScale == 1.0 {
		var idx int = 0
		for row := 0; row < int(hist.height); row++ {
			histRow := rs.rowVals(row)
			for x, val := range histRow {
				sscale := supScale(x, row, centx, centy, hist.grad.MaxMod)
				suppressed[idx] = float64(val) * sscale
				if suppressed[idx] > maxSuppressed {
					maxSuppressed = suppressed[idx]
				}
				idx++
			}
		}
		var pixScale float64 = 255.0 / maxSuppressed
		//fmt.Println("Suppressed Render pixScale factor:", pixScale)
		rndPix := rnd.Pix()
		for index, val := range suppressed {
			rndPix[index] = uint8(val * pixScale)
		}
	} else {
		// Scale while rendering
	}
	return rnd
}

// setupInvertedBins populates the invertedBins map for the given histogram.
// invertedBins stores the index in the bins slice for each occurring histogram
// value.
// Used to lazily populate invertedBins when needed.
func setupInvertedBins(bins []BinPair) (invertedBins map[uint32]int) {
	invertedBins = make(map[uint32]int, len(bins))
	for index, val := range bins {
		invertedBins[val.BinVal] = index
	}
	return
}

// RenderSubstituteCore renders an 8-bit image of the histogram, substituting
// the given value as the pixel value for each corresponding bin value. The
// input slice must be the same length as the slice of bin values returned
// by Bins, and contain new values corresponding to that order.
// This is used to render the delentropy values of the histogram.
// Note that as the slice returned by Bins() does not include 0 values,
// the value to be used for empty bins must be supplied.
func (hist *histCore) renderSubstituteCore(rs rowSource, subs []uint8, zeroVal uint8) SippImage {
	if hist.invertedBins == nil {
		hist.invertedBins = setupInvertedBins(hist.bins)
	}
	rnd, scale, _ := hist.renderInto()
	rndPix := rnd.Pix()
	if scale == 1.0  {
		for row := 0; row < hist.height; row++ {
			histRow := rs.rowVals(row)
			for index, val := range histRow {
				if val == 0 {
					rndPix[row*hist.width+index] = zeroVal;
				} else {
					rndPix[row*hist.width+index]= subs[hist.invertedBins[val]]
				}
			}
		}
	} else {
		// Scale while rendering
	}
	return rnd
}
