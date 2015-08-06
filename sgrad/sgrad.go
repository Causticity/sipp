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
	Pix []complex64
	Rect image.Rectangle
	MaxMod float64
}

// Use a 2x2 kernel to create a finite-differences gradient image, one pixel
// narrower and shorter than the original. We'd rather reduce the size of the
// output image than arbitrarily wrap around or extend the source image, as
// any such procedure could introduce errors into the statistics.
func Fdgrad(src *Sippimage) (grad *Gradimage) {
	// Create the dst image from the bounds of the src
	srect := src.Img.Bounds()
	grad = new(Gradimage)
	grad.Rect = image.Rect(0,0,srect.Dx()-1,srect.Dy()-1)
	grad.Pix = make([]complex64, grad.Rect.Dx()*grad.Rect.Dy())
	grad.MaxMod = 0
	
	fmt.Println("source image rect:<", srect, ">")
	fmt.Println("source image stride:", src.Img.Stride)
	fmt.Println("gradient image rect:<", grad.Rect, ">")
	fmt.Println("Gradient image no. of pixels:<", len(grad.Pix), ">")
	
	// Drive over the dst image
	// grad[x,y] = complex(src[x+1,y+1] - src[x,y], src[x+1,y]-src[x,y+1])
	// loop over dest scanlines
    gradStride := grad.Rect.Dx()
	for line := 0; line < grad.Rect.Dy(); line++ {
		// Set the following slice indices into the src:
		cur := src.Img.PixOffset(0, line)
		rightdown := src.Img.PixOffset(1, line+1)
		right := src.Img.PixOffset(1, line)
		down := src.Img.PixOffset(0, line+1)
		dstMin := line*gradStride
		dstMax := dstMin+gradStride
		for dsti := dstMin; dsti < dstMax; dsti++ {
			re := float32(src.Img.Pix[rightdown]) - float32(src.Img.Pix[cur])
			im := float32(src.Img.Pix[right]) - float32(src.Img.Pix[down])
			grad.Pix[dsti] = complex(re, im)
			modsq := float64(re)*float64(re) + float64(im)*float64(im)
			if modsq > grad.MaxMod {
				grad.MaxMod = modsq
			}
			cur++
			rightdown++
			right++
			down++
		}
	}
	grad.MaxMod = math.Sqrt(grad.MaxMod)

	return
}

func (grad *Gradimage) Render() (re, im *Sippimage) {
	// compute max excursions of the real and imag parts
	var minreal float32 = math.MaxFloat32
	var minimag float32 = math.MaxFloat32
	var maxreal float32 = -math.MaxFloat32
	var maximag float32 = -math.MaxFloat32
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
	re = new(Sippimage)
	re.Img = image.NewGray(grad.Rect)
	im = new(Sippimage)
	im.Img = image.NewGray(grad.Rect)
	// scan the complex image, generating the two renderings
	for index, pix := range grad.Pix {
		r := real(pix)
		i := imag(pix)
		re.Img.Pix[index] = uint8((r - minreal)*rscale)
		im.Img.Pix[index] = uint8((i - minimag)*iscale)
	}
	return
}