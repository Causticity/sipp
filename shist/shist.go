package shist

import (
	"fmt"
	"math"
)

import (
	. "github.com/Causticity/sipp/sgrad"
	. "github.com/Causticity/sipp/simage"
)

type Sipphist struct {
	grad *Gradimage
	bin [] uint32
	max uint32
	k int
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
		u := math.Floor(factor*float64(real(pixel)) + float64(k)) 
		v := math.Floor(factor*float64(imag(pixel)) + float64(k))
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

func (hist *Sipphist) Render() (rnd *Sippimage) {
	return
}