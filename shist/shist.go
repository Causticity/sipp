// Copyright Raul Vera 2015-2016

// Package shist provides functions for computing and rendering a 2-dimensional
// histogram of values of a complex gradient image, as well as a 2D Entropy
// calculation based on the histogram.
// TODO: This now has more than just the 2D histogram: 
//    - 1-D
//    - both kinds of entropy
//   Refactor this along cleaner lines. Some to sgrad.go, some to an entropy
//   package.
package shist

import (
	"fmt"
	"image"
	"math"
)

import (
	. "github.com/Causticity/sipp/sgrad"
	. "github.com/Causticity/sipp/simage"
)

// SippHist is a 2-dimensional histogram of the values in a complex gradient
// image.
type SippHist struct {
	// A reference to the gradient image we are computing from
	grad *GradImage
	// The size of our histogram. It will be 2*k+1 on a side.
	k uint16
	// The histogram data.
	bin [] uint32
	// The index of the histogram bin for each gradient image pixel.
	binIndex [] int
	// The maximum bin value in the histogram.
	max uint32
	// A suppressed version of the histogram, stored as floats for subsequent
	// computation.
	suppressed [] float64
	// The maximum suppressed value, stored as a float for subsequent 
	// computation.
	maxSuppressed float64
	// The delentropy for each bin value that actually occurred.
	delentropy [] float64
	// The largest delentropy value.
	maxDelentropy float64
}

const maxK = 2048
const kMargin = 8

const histSize8BPP = 256
const histSize16BPP = 65536

// GreyHist computes a 1D histogram of the greyscale values in the image.
func GreyHist(im SippImage) ([]uint32) {
	histSize := histSize8BPP
	is16 := false
	if im.Bpp() == 16 {
		histSize = histSize16BPP
		is16 = true
	}
	
	//fmt.Println("GreyHist histogram size is ", histSize)
	
	hist := make([] uint32, histSize)
	imPix := im.Pix()
	for y := 0; y < im.Bounds().Dy(); y++ {
		for x:= 0; x < im.Bounds().Dx(); x++ {
			index := im.PixOffset(x, y)
			var val uint16 = uint16(imPix[index])
			if is16 {
				val = val << 8 | uint16(imPix[index+1])
			}
			hist[val]++
		}
	}
	return hist
}

// Entropy calculates the conventional entropy of the image.
func Entropy(im SippImage) (float64, SippImage) {
	hist := GreyHist(im)
	total := float64(im.Bounds().Dx()*im.Bounds().Dy())
	normHist := make([]float64, len(hist))
	var check float64
	for i, binVal := range hist {
		normHist[i] = float64(binVal)/total
		check = check + normHist[i]
	}
	//fmt.Println("Normalised histogram sums to ", check)
	entHist := make ([]float64, len(hist))
	var ent, maxEnt float64
	for j, p := range normHist {
		if p > 0 {
			entHist[j] = -1.0 * p * math.Log2(p)
			ent = ent + entHist[j]
			if entHist[j] > maxEnt {
				maxEnt = entHist[j]
			}
		}
	}
	//fmt.Println("maxEnt is ", maxEnt)
	entIm := new(SippGray)
	entIm.Gray = image.NewGray(im.Bounds())
	entImPix := entIm.Pix()

	// scale the entropy from (0-maxEnt) to (0-255)
	is16 := false
	if im.Bpp() == 16 {
		is16 = true
	}
	scale := 255.0 / maxEnt
	width := im.Bounds().Dx()
	imPix := im.Pix()
	for y := 0; y < im.Bounds().Dy(); y++ {
		for x := 0; x < width; x++ {
			index := im.PixOffset(x, y)
			var val uint16 = uint16(imPix[index])
			if is16 {
				val = val << 8 | uint16(imPix[index+1])
			}
			entImPix[y*width+x] = uint8(math.Floor(entHist[val] * scale))
		}
	}
	return ent, entIm
}

// Hist computes the 2D histogram, 2*K=1 on a side with 0,0 at the center, from
// the given gradient image.
// TODO: For 16-bit images, this should use sparse techniques, because the max
// mod may be huge, but there will only ever be as many values as pixels. In
// general, we should avoid scaling down. We should always be able to deal with
// the exact histogram, without any scaling into the bins. So perhaps get rid of
// "k"? But the number of bins does affect the entropy, 
func Hist(grad *GradImage, k uint16) (hist *SippHist) {
	if k == 0 {
		if grad.MaxMod > maxK {
			k = maxK
			fmt.Println("clamping k; max mod is ", grad.MaxMod)
		} else {
			k = uint16(grad.MaxMod)+kMargin
		}
	}
	fmt.Println("K:", k, " histogram edge size:", (k*2+1))
	// create the 2D histogram bins as 2K+1 on a side, so always odd
	hist = new(SippHist)
	hist.grad = grad
	hist.k = k
	hist.binIndex = make([] int, len(grad.Pix))
	stride := int(2*k+1)
	histDataSize := stride*stride
	hist.bin = make([] uint32, histDataSize)
	//fmt.Println("histogram side:", stride, " data size: ", histDataSize)

	// Walk through the image, computing the bin address from the gradient 
	// values storing the bin address in binIndex and incrementing the bin.
	// Save the maximum bin value as well.
	var factor float64 = 1.0
	if grad.MaxMod > float64(k) {
	    factor = float64(k) / grad.MaxMod
	}
	//fmt.Println("MaxMod:", grad.MaxMod, " factor:", factor)
	for i, pixel := range grad.Pix {
		u := int(math.Floor(factor*real(pixel))) + int(k)
		v := int(math.Floor(factor*imag(pixel))) + int(k)
		hist.binIndex[i] = v*stride + u
		hist.bin[hist.binIndex[i]]++
		if hist.bin[hist.binIndex[i]] > hist.max {
			hist.max = hist.bin[hist.binIndex[i]]
		}
	}
	//fmt.Println("Histogram complete. Maximum bin value:", hist.max)
	return
}

// Delentropy returns the 2D entropy of the gradient image.
func (hist *SippHist) Delentropy() (float64) {
	// Store the entropy values corresponding to the bin counts that actually
	// occurred.
	hist.delentropy = make([] float64, hist.max+1)
    total := float64(len(hist.grad.Pix))
    hist.maxDelentropy = 0.0
    var dent float64 = 0.0
	for _, bin := range hist.bin {
		if bin != 0 {
			// compute the entropy only once for a given bin value.
			if hist.delentropy[bin] == 0.0 {
				p := float64(bin) / total
				hist.delentropy[bin] = p * math.Log2(p) * -1.0
			}
			dent += hist.delentropy[bin]
			if hist.delentropy[bin] > hist.maxDelentropy {
				hist.maxDelentropy = hist.delentropy[bin]
			}
		}
	}
	
	return dent
}

// HistDelentropyImage returns a greyscale image of the delentropy for each
// histogram bin. DelEntropy must have been called first.
func (hist *SippHist) HistDelentropyImage() (SippImage) {
	// Ensure that we have the table of delentropies and the maximum
	if hist.delentropy == nil {
		fmt.Println("Warning: HistDelntropyImage called before computing delentropy!")
		return nil
	}
	// Make a greyscale image of the entropy for each bin.
	stride := int(2*hist.k+1)
	dentGray := new(SippGray)
	dentGray.Gray = image.NewGray(image.Rect(0,0,stride,stride))
	dentGrayPix := dentGray.Pix()
	// scale the delentropy from (0-hist.maxDelentropy) to (0-255)
	scale := 255.0 / hist.maxDelentropy
	for i, val := range hist.bin {
		dentGrayPix[i] = uint8(hist.delentropy[val]*scale)
	}
	return dentGray
}

// DelEntropyImage returns a greyscale image of the entropy for each gradient
// pixel. DelEntropy must have been called first.
func (hist *SippHist) DelEntropyImage() (SippImage) {
	// Ensure that we have the table of entropies and the maximum
	if hist.delentropy == nil {
		fmt.Println("Warning: DelEntropyImage called before computing delentropy!")
		return nil
	}
	// Make a greyscale image of the entropy for each bin.
	dentGray := new(SippGray)
	dentGray.Gray = image.NewGray(hist.grad.Rect)
	dentGrayPix := dentGray.Pix()
	// scale the entropy from (0-hist.maxEntropy) to (0-255)
	scale := 255.0 / hist.maxDelentropy
	for i := range dentGrayPix {
		// The following lookup works as follows:
		// i - the gradient (and delentropy) image-pixel index
		// hist.binIndex[i] - the histogram bin for that pixel
		// hist.bin[hist.binIndex[i] - the value in that bin
		// hist.delentropy[...] The delentropy for that value
		// We scale that delentropy and convert to an 8-bit pixel
		dentGrayPix[i] = uint8(hist.delentropy[hist.bin[hist.binIndex[i]]]*scale)
	}
	return dentGray
}

func supScale(x, y, k int) float64 {
	xdist := float64(x - k)
	ydist := float64(y - k)
	hyp := math.Hypot(xdist, ydist)
	return (hyp/float64(k))
}

// Suppress suppresses the spike near the origin of the histogram by scaling the
// values in the histogram by their distance from the origin. All the values in
// in the histogram are multiplied by a factor that ranges from 0.0 at the
// origin to 1.0 at k in any direction from the origin.
func (hist *SippHist) Suppress() {
	if hist.suppressed != nil {
		return
	}
	stride := int(2*hist.k+1)
	size := stride*stride
	hist.suppressed = make([]float64, size)
	var index uint32 = 0
	hist.maxSuppressed = 0
	for y := 0; y < stride; y++ {
		for x := 0; x < stride; x++ {
			sscale := supScale(x, y, int(hist.k))
			hist.suppressed[index] = float64(hist.bin[index]) * sscale
			if hist.suppressed[index] > hist.maxSuppressed {
				hist.maxSuppressed = hist.suppressed[index]
			}
			index++
		}
	}
	//fmt.Println("Distance suppression complete; max suppressed value:", hist.maxSuppressed)
}

// RenderSuppressed renders the suppressed version of the histogram and returns
// the result as an 8-bit grayscale image.
func (hist *SippHist) RenderSuppressed() SippImage {
	// Here we will generate an 8-bit output image of the same size as the
	// histogram, scaled to use the full dynamic range of the image format.
	hist.Suppress()
	stride := int(2*hist.k+1)
	var scale float64 = 255.0 / hist.maxSuppressed
	//fmt.Println("Suppressed Render scale factor:", scale)
	rnd := new(SippGray)
	rnd.Gray = image.NewGray(image.Rect(0,0,stride,stride))
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
	stride := int(2*hist.k+1)
	//var scale float64 = 255.0 / float64(hist.max)
	//fmt.Println("Render scale factor:", scale)
	rnd := new(SippGray)
	rnd.Gray = image.NewGray(image.Rect(0,0,stride,stride))
	rndPix := rnd.Pix()
	for index, val := range hist.bin {
		if val > 255 {
			val = 255
		}
		rndPix[index] = uint8(val)
	}
	return rnd
}
