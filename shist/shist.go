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
	Grad *Gradimage
	Bin [] uint32
	Max uint32
	K int
}

func Hist(grad *Gradimage, k int) (hist *Sipphist) {
	// create the 2D histogram bins as 2K+1 on a side, so always odd
	hist = new(Sipphist)
	hist.Grad = grad
	hist.K = k
	stride := 2*k+1
	histDataSize := stride*stride
	hist.Bin = make([] uint32, histDataSize)
	fmt.Println("histogram side:", stride, " data size: ", histDataSize)
	
	// Walk through the image, computing the bin address from the gradient 
	// values. 
    gradStride := grad.Rect.Dx()
    factor := float64(k) / grad.MaxMod
    fmt.Println("MaxMod:", grad.MaxMod, " factor:", factor) 
	for line := 0; line < grad.Rect.Dy(); line++ {
		lineMin := line*gradStride
		lineMax := lineMin+gradStride
		for gradIndex := lineMin; gradIndex < lineMax; gradIndex++ {
			u := math.Floor(factor*float64(real(grad.Pix[gradIndex])) + float64(k)) 
			v := math.Floor(factor*float64(imag(grad.Pix[gradIndex])) + float64(k))
			histIndex := int(v)*stride + int(u)
			if histIndex <0 || histIndex >= histDataSize {
				fmt.Println("YOW! index out of range at line ", line, ":", histIndex,
							" u:", u, " v:", v, "grad.Pix[gradIndex]:", grad.Pix[gradIndex])
				return
			}
			hist.Bin[histIndex]++
			if hist.Bin[histIndex] > hist.Max {
				hist.Max = hist.Bin[histIndex]
			}
		}		
	}
	fmt.Println("Histogram complete. Maximum bin value:", hist.Max)
	return
}

func (hist *Sipphist) Render() (rnd *Sippimage) {
	return
}