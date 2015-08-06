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

type Sipphist struct {
	grad *Gradimage
	k int
	bin [] uint32
	max uint32
	suppressed [] float64
	maxSuppressed float64
}

func Hist(grad *Gradimage, k int) (hist *Sipphist) {
	// create the 2D histogram bins as 2K+1 on a side, so always odd
	hist = new(Sipphist)
	hist.grad = grad
	hist.k = k
	stride := 2*k+1
	histDataSize := stride*stride
	hist.bin = make([] uint32, histDataSize)
	fmt.Println("histogram side:", stride, " data size: ", histDataSize)
	
	// Walk through the image, computing the bin address from the gradient 
	// values. 
    factor := float64(k) / grad.MaxMod
    fmt.Println("MaxMod:", grad.MaxMod, " factor:", factor)
	for _, pixel := range grad.Pix {
		u := factor*real(pixel) + float64(k)
		v := factor*imag(pixel) + float64(k)
		histIndex := int(v)*stride + int(u)
		hist.bin[histIndex]++
		if hist.bin[histIndex] > hist.max {
			hist.max = hist.bin[histIndex]
		}
	}
	fmt.Println("Histogram complete. Maximum bin value:", hist.max)
	return
}

func (hist *Sipphist) Entropy() (ent float64) {
    total := float64(len(hist.grad.Pix))
	for _, bin := range hist.bin {
		if bin != 0 {
			p := float64(bin) / total
			ent += p * math.Log2(p)
		}
	}
	ent *= -1
	return
}

func supScale(x, y, k int) float64 {
	xdist := float64(x - k)
	ydist := float64(y - k)
	hyp := math.Hypot(xdist, ydist)
	return (hyp/float64(k))
}

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

func (hist *Sipphist) RenderSuppressed() (rnd *Sippimage) {
	// Here we will generate an 8-bit output image of the same size as the
	// histogram, scaled to use the full dynamic range of the image format.
	hist.Suppress()
	stride := 2*hist.k+1
	var scale float64 = 255.0 / hist.maxSuppressed
	fmt.Println("Suppressed Render scale factor:", scale)
	rnd = new(Sippimage)
	rnd.Img = image.NewGray(image.Rect(0,0,stride,stride))
	for index, val := range hist.suppressed {
		rnd.Img.Pix[index] = uint8(val * scale)
	}
	return
}

func (hist *Sipphist) Render() (rnd *Sippimage) {
	// Here we will generate an 8-bit output image of the same size as the
	// histogram, clipped to 255.
	stride := 2*hist.k+1
	//var scale float64 = 255.0 / float64(hist.max)
	//fmt.Println("Render scale factor:", scale)
	rnd = new(Sippimage)
	rnd.Img = image.NewGray(image.Rect(0,0,stride,stride))
	for index, val := range hist.bin {
		//rnd.Img.Pix[index] = uint8(float64(val) * scale)
		if val > 255 {
			val = 255
		}
		rnd.Img.Pix[index] = uint8(val)
	}
	return
}