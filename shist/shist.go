// Package shist provides functions for computing and rendering a 2-dimensional
// histogram of values of a complex gradient image, as well as a 2D Entropy
// calculation based on the histogram.
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

// Sipphist is a 2-dimensional histogram of the values in a complex gradient
// image.
type Sipphist struct {
	// A reference to the gradient image we are computing from
	grad *Gradimage
	// The size of our histogram. It will be 2*k+1 on a side.
	k int
	// The histogram data.
	bin [] uint32
	// The maximum bin value in the histogram.
	max uint32
	// A suppressed version of the histogram, stored as floats for subsequent
	// computation.
	suppressed [] float64
	// The maximum suppressed value, stored as a float for subsequent 
	// computation.
	maxSuppressed float64
	// The entropy for each bin value that actually occurred.
	entropy [] float64
	// The largest entropy value.
	maxEntropy float64
}

const maxK = 2048
const kMargin = 8

// Hist computes the 2D histogram, 2*K=1 on a side with 0,0 at the center, from
// the given gradient image.
func Hist(grad *Gradimage, k int) (hist *Sipphist) {
	if k == 0 {
		if grad.MaxMod > maxK {
			k = maxK
		} else {
			k = int(grad.MaxMod)+kMargin
		}
	}
	fmt.Println("K:", k, " histogram edge size:", (k*2+1))
	// create the 2D histogram bins as 2K+1 on a side, so always odd
	hist = new(Sipphist)
	hist.grad = grad
	hist.k = k
	stride := 2*k+1
	histDataSize := stride*stride
	hist.bin = make([] uint32, histDataSize)
	fmt.Println("histogram side:", stride, " data size: ", histDataSize)
	
	// Walk through the image, computing the bin address from the gradient 
	// values and incrementing the bin.
	// Save the maximum bin value as well.
	var factor float64 = 1.0
	if grad.MaxMod > float64(k) {
	    factor = float64(k) / grad.MaxMod
	}
	fmt.Println("MaxMod:", grad.MaxMod, " factor:", factor)
	for _, pixel := range grad.Pix {
		u := int(math.Floor(factor*real(pixel))) + k
		v := int(math.Floor(factor*imag(pixel))) + k
		binIndex := v*stride + u
		hist.bin[binIndex]++
		if hist.bin[binIndex] > hist.max {
			hist.max = hist.bin[binIndex]
		}
	}
	fmt.Println("Histogram complete. Maximum bin value:", hist.max)
	return
}


// Entropy returns the 2D entropy of the gradient image, and a greyscale image
// of the entropy for each histogram bin.
func (hist *Sipphist) Entropy() (float64, SippImage) {
	// Store the entropy values corresponding to the bin counts that actually
	// occurred.
	hist.entropy = make([] float64, hist.max+1)
    total := float64(len(hist.grad.Pix))
    hist.maxEntropy = 0.0
    var ent float64 = 0.0
	for _, bin := range hist.bin {
		if bin != 0 {
			// compute the entropy only once for a given bin value.
			if hist.entropy[bin] == 0.0 {
				p := float64(bin) / total
				hist.entropy[bin] = p * math.Log2(p) * -1.0
			}
			ent += hist.entropy[bin]
			if hist.entropy[bin] > hist.maxEntropy {
				hist.maxEntropy = hist.entropy[bin]
			}
		}
	}
	// Now that we have the table of entropies and the maximum, make a 
	// greyscale image of the entropy for each bin.
	stride := 2*hist.k+1
	entGray := new(SippGray)
	entGray.Gray = image.NewGray(image.Rect(0,0,stride,stride))
	entGrayPix := entGray.Pix()
	// scale the entropy from (0-hist.maxEntropy) to (0-255)
	scale := 255.0 / hist.maxEntropy
	for i, val := range hist.bin {
		entGrayPix[i] = uint8(hist.entropy[val]*scale)
	}
	return ent, entGray
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
func (hist *Sipphist) Suppress() {
	if hist.suppressed != nil {
		return
	}
	stride := 2*hist.k+1
	size := stride*stride
	hist.suppressed = make([]float64, size)
	var index uint32 = 0
	hist.maxSuppressed = 0
	for y := 0; y < stride; y++ {
		for x := 0; x < stride; x++ {
			sscale := supScale(x, y, hist.k)
			hist.suppressed[index] = float64(hist.bin[index]) * sscale
			if hist.suppressed[index] > hist.maxSuppressed {
				hist.maxSuppressed = hist.suppressed[index]
			}
			index++
		}
	}
	fmt.Println("Distance suppression complete; max suppressed value:", hist.maxSuppressed)
}

// RenderSuppressed renders the suppressed version of the histogram and returns
// the result as an 8-bit grayscale image.
func (hist *Sipphist) RenderSuppressed() SippImage {
	// Here we will generate an 8-bit output image of the same size as the
	// histogram, scaled to use the full dynamic range of the image format.
	hist.Suppress()
	stride := 2*hist.k+1
	var scale float64 = 255.0 / hist.maxSuppressed
	fmt.Println("Suppressed Render scale factor:", scale)
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
func (hist *Sipphist) Render() SippImage {
	// Here we will generate an 8-bit output image of the same size as the
	// histogram, clipped to 255.
	stride := 2*hist.k+1
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