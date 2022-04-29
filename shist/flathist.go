// Copyright Raul Vera 2015-2021

package shist

import (
	//"fmt"
	"image"
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
		var found bool
		if binval != 0 {
			found = false
			//fmt.Printf("looking for binval %d\n", binval)
			for i, pair := range hist.bins {
				if binval == pair.BinVal {
					//fmt.Println("Found it, incrementing pair.num")
					hist.bins[i].Num++
					found = true
				}
			}
			if found == false {
				//fmt.Printf("Appending pair for binval %d\n", binval)
				hist.bins = append(hist.bins, BinPair{binval, 1})
			}
		}
	}
	//fmt.Println("Histogram complete. Maximum bin value:", hist.Max)
	return hist
}

// Bins returns a compact slice of the actually used bins. There is no order
// specified, but each call to Bins returns the values in the same order.
func (hist *flatSippHist) Bins() ([]BinPair) {
	return hist.bins
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

// suppress returns a suppressed version of the bins of the histogram, along
// with the maximum suppressed value.
func (hist *flatSippHist) suppress() (suppressed []float64, maxSuppressed float64 ){
	centx := (int(hist.width)-1)/2
	centy := (int(hist.height)-1)/2

	size := int(hist.width) * int(hist.height)
	suppressed = make([]float64, size)
	var index int = 0
	for y := 0; y < int(hist.height); y++ {
		for x := 0; x < int(hist.width); x++ {
			sscale := supScale(x, y, centx, centy, hist.grad.MaxMod)
			suppressed[index] = float64(hist.bin[index]) * sscale
			if suppressed[index] > maxSuppressed {
				maxSuppressed = suppressed[index]
			}
			index++
		}
	}
	//fmt.Println("Distance suppression complete; max suppressed value:", supp.Max)
	return
}

// Render renders the histogram into an 8-bit grayscale image. If clip is true,
// values are clipped to 255. If clip is false, values are scaled to 255.
func (hist *flatSippHist) Render(clip bool) SippImage {
	width, height := hist.Size()
	rnd := new(SippGray)
	rnd.Gray = image.NewGray(image.Rect(0, 0, int(width), int(height)))
	rndPix := rnd.Pix()
	var scale float64 = 255.0 / float64(hist.max)
	//fmt.Println("Render scale factor:", scale)
	for index, val := range hist.bin {
		if clip {
			if val > 255 {
				val = 255
			}
		} else {
			val = uint32(math.Round(float64(val) * scale))
		}
		rndPix[index] = uint8(val)
	}
	return rnd
}

// RenderSuppressed renders a suppressed version of the histogram and returns
// the result as an 8-bit grayscale image.
func (hist *flatSippHist) RenderSuppressed() SippImage {
	suppressed, maxSuppressed := hist.suppress()
	width, height := hist.Size()
	var scale float64 = 255.0 / maxSuppressed
	//fmt.Println("Suppressed Render scale factor:", scale)
	rnd := new(SippGray)
	rnd.Gray = image.NewGray(image.Rect(0, 0, int(width), int(height)))
	rndPix := rnd.Pix()
	for index, val := range suppressed {
		rndPix[index] = uint8(val * scale)
	}
	return rnd
}

// setupInvertedBins populates the invertedBins map for the given histogram.
// invertedBins stores the index in the bins slice for each occurring histogram
// value.
// Used to lazily populate invertedBins when needed.
func setupInvertedBins(hist *flatSippHist) {
	if hist.invertedBins != nil {
		return;
	}
	hist.invertedBins = make(map[uint32]int, len(hist.bins))
	for index, val := range hist.bins {
		hist.invertedBins[val.BinVal] = index
	}
}

// RenderSubstitute renders an 8-bit image of the histogram, substituting
// the given value as the pixel value for each corresponding bin value. The
// input slice must be the same length as the slice of bin values returned
// by Bins, and contain new values corresponding to that order.
// This is used to render the delentropy values of the histogram.
// Note that as the slice returned by Bins() does not include 0 values,
// the value to be used for empty bins must be supplied.
func (hist *flatSippHist) RenderSubstitute(subs []uint8, zeroVal uint8) SippImage {
	if hist.invertedBins == nil {
		setupInvertedBins(hist)
	}
	width, height := hist.Size()
	rnd := new(SippGray)
	rnd.Gray = image.NewGray(image.Rect(0, 0, int(width), int(height)))
	rndPix := rnd.Pix()
	for index, val := range hist.bin {
		if val == 0 {
			rndPix[index] = zeroVal;
		} else {
			rndPix[index]= subs[hist.invertedBins[val]]
		}
	}
	return rnd
}
