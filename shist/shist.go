// Copyright Raul Vera 2015-2020

// Package shist provides functions for computing a histogram of values of an
// image, and for computing and rendering a 2-dimensional histogram of values of
// a complex gradient image.
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
	// The size of our histogram. It will be 2*Radius+1 on a side.
	Radius uint16
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

const maxRadius = 2048
const radiusMargin = 8

const histSize8BPP = 256
const histSize16BPP = 65536

// GreyHist computes a 1D histogram of the greyscale values in the image.
func GreyHist(im SippImage) []uint32 {
	histSize := histSize8BPP
	is16 := false
	if im.Bpp() == 16 {
		histSize = histSize16BPP
		is16 = true
	}

	//fmt.Println("GreyHist histogram size is ", histSize)

	hist := make([]uint32, histSize)
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
	return hist
}

// Hist computes the 2D histogram, 2*radius+1 on a side with 0,0 at the center,
// from the given gradient image.
// TODO: For 16-bit images, this should use sparse techniques, because the max
// mod may be huge, but there will only ever be as many values as pixels. In
// general, we should avoid scaling down. We should always be able to deal with
// the exact histogram, without any scaling into the bins. So perhaps get rid of
// the "radius" argument entirely?
func Hist(grad *ComplexImage, radius uint16) (hist *SippHist) {
	if radius == 0 {
		if grad.MaxMod > maxRadius {
			radius = maxRadius
			fmt.Println("clamping radius; max mod is ", grad.MaxMod)
		} else {
			radius = uint16(grad.MaxMod) + radiusMargin
		}
	}
	fmt.Println("Radius:", radius, " histogram edge size:", (radius*2 + 1))
	// create the 2D histogram bins as 2radius+1 on a side, so always odd
	hist = new(SippHist)
	hist.Grad = grad
	hist.Radius = radius
	hist.BinIndex = make([]int, len(grad.Pix))
	stride := int(2*radius + 1)
	histDataSize := stride * stride
	hist.Bin = make([]uint32, histDataSize)
	//fmt.Println("histogram side:", stride, " data size: ", histDataSize)

	// Walk through the image, computing the bin address from the gradient
	// values storing the bin address in BinIndex and incrementing the bin.
	// Save the maximum bin value as well.
	var factor float64 = 1.0
	if grad.MaxMod > float64(radius) {
		factor = float64(radius) / grad.MaxMod
	}
	//fmt.Println("MaxMod:", grad.MaxMod, " factor:", factor)
	for i, pixel := range grad.Pix {
		u := int(math.Floor(factor*real(pixel))) + int(radius)
		v := int(math.Floor(factor*imag(pixel))) + int(radius)
		hist.BinIndex[i] = v*stride + u
		hist.Bin[hist.BinIndex[i]]++
		if hist.Bin[hist.BinIndex[i]] > hist.Max {
			hist.Max = hist.Bin[hist.BinIndex[i]]
		}
	}
	//fmt.Println("Histogram complete. Maximum bin value:", hist.Max)
	return
}

func supScale(x, y, radius int) float64 {
	xdist := float64(x - radius)
	ydist := float64(y - radius)
	hyp := math.Hypot(xdist, ydist)
	return (hyp / float64(radius))
}

// Suppress suppresses the spike near the origin of the histogram by scaling the
// values in the histogram by their distance from the origin. All the values in
// in the histogram are multiplied by a factor that ranges from 0.0 at the
// origin to 1.0 at radius in any direction from the origin.
func (hist *SippHist) Suppress() {
	if hist.suppressed != nil {
		return
	}
	stride := int(2*hist.Radius + 1)
	size := stride * stride
	hist.suppressed = make([]float64, size)
	var index uint32 = 0
	hist.maxSuppressed = 0
	for y := 0; y < stride; y++ {
		for x := 0; x < stride; x++ {
			sscale := supScale(x, y, int(hist.Radius))
			hist.suppressed[index] = float64(hist.Bin[index]) * sscale
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
	stride := int(2*hist.Radius + 1)
	var scale float64 = 255.0 / hist.maxSuppressed
	//fmt.Println("Suppressed Render scale factor:", scale)
	rnd := new(SippGray)
	rnd.Gray = image.NewGray(image.Rect(0, 0, stride, stride))
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
	stride := int(2*hist.Radius + 1)
	//var scale float64 = 255.0 / float64(hist.Max)
	//fmt.Println("Render scale factor:", scale)
	rnd := new(SippGray)
	rnd.Gray = image.NewGray(image.Rect(0, 0, stride, stride))
	rndPix := rnd.Pix()
	for index, val := range hist.Bin {
		if val > 255 {
			val = 255
		}
		rndPix[index] = uint8(val)
	}
	return rnd
}
