package sgrad

// Create a gradient image from a source image.

import (
	"image"
    "fmt"
    "math"
)

import (
	. "github.com/Causticity/sipp/simage"
)

type Gradimage struct {
	Pix []complex128
	Rect image.Rectangle
	MaxMod float64
}

// Use a 2x2 kernel to create a finite-differences gradient image, one pixel
// narrower and shorter than the original. We'd rather reduce the size of the
// output image than arbitrarily wrap around or extend the source image, as
// any such procedure could introduce errors into the statistics.
func Fdgrad(src Sippimage) (grad *Gradimage) {
	// Create the dst image from the bounds of the src
	srect := src.Bounds()
	grad = new(Gradimage)
	grad.Rect = image.Rect(0,0,srect.Dx()-1,srect.Dy()-1)
	grad.Pix = make([]complex128, grad.Rect.Dx()*grad.Rect.Dy())
	grad.MaxMod = 0
	
	fmt.Println("source image rect:<", srect, ">")
	fmt.Println("source image stride:", src.Stride())
	fmt.Println("gradient image rect:<", grad.Rect, ">")
	fmt.Println("Gradient image no. of pixels:<", len(grad.Pix), ">")
	
	// Drive over the dst image
	// grad[x,y] = complex(src[x+1,y+1] - src[x,y], src[x+1,y]-src[x,y+1])
	// loop over dest scanlines
    gradStride := grad.Rect.Dx()
    pix := src.Pix()
    step := src.Step()
	for line := 0; line < grad.Rect.Dy(); line++ {
		// Set the following slice indices into the src:
		cur := src.PixOffset(0, line)
		rightdown := src.PixOffset(1, line+1)
		right := src.PixOffset(1, line)
		down := src.PixOffset(0, line+1)
		dstMin := line*gradStride
		dstMax := dstMin+gradStride
		for dsti := dstMin; dsti < dstMax; dsti++ {
			// This is not polymorphic access. Add a method that returns the
			// pixel value as a float64 and just use the indices
			re := float64(pix[rightdown]) - float64(pix[cur])
			im := float64(pix[right]) - float64(pix[down])
			grad.Pix[dsti] = complex(re, im)
			modsq := re*re + im*im
			if modsq > grad.MaxMod {
				grad.MaxMod = modsq
			}
			cur += step
			rightdown += step
			right += step
			down += step
		}
	}
	grad.MaxMod = math.Sqrt(grad.MaxMod)

	return
}

func (grad *Gradimage) Render() (Sippimage, Sippimage) {
	// compute max excursions of the real and imag parts
	var minreal float64 = math.MaxFloat64
	var minimag float64 = math.MaxFloat64
	var maxreal float64 = -math.MaxFloat64
	var maximag float64 = -math.MaxFloat64
	for _, pix := range grad.Pix {
		re := real(pix)
		im := imag(pix)
		if re < minreal {
			minreal = re
		}
		if re > maxreal {
			maxreal = re
		}
		if im < minimag {
			minimag = im
		}
		if im > maximag {
			maximag = im
		}
	}
	// compute scale and offset for each image
	rscale := 255.0 / (maxreal - minreal)
	iscale := 255.0 / (maximag - minimag)
	re := new(SippGray)
	re.Gray = image.NewGray(grad.Rect)
	im := new(SippGray)
	im.Gray = image.NewGray(grad.Rect)
	// scan the complex image, generating the two renderings
	rePix := re.Pix()
	imPix := im.Pix()
	for index, pix := range grad.Pix {
		r := real(pix)
		i := imag(pix)
		rePix[index] = uint8((r - minreal)*rscale)
		imPix[index] = uint8((i - minimag)*iscale)
	}
	return re, im
}