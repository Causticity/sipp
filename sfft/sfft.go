// Computes an FFT using fftw

package sfft

import (
	"image"
//   "fmt"
    "math"

    "github.com/runningwild/go-fftw/fftw"
)

import (	
	. "github.com/Causticity/sipp/simage"
	. "github.com/Causticity/sipp/scomplex"
)

type FFTimage struct {
	Compleximage
}

func FFT(src SippImage) (fft *FFTimage) {
	fft = new(FFTimage)
	fft.Rect = src.Bounds()
	width := fft.Rect.Dx()
	height := fft.Rect.Dy()
	size := width*height
	fft.Pix = make([]complex128, size)
	// Multiply by (-1)^(x+y) while converting the pixels to complex numbers
	shiftStart := 1.0
	shift := shiftStart
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fft.Pix[i] = complex(src.Val(x,y)*shift, 0)
			i++
			shift = -shift
		}
		shiftStart = -shiftStart
		shift = shiftStart
	}

	re, im := fft.Render()
	reName := "prefftreal.png"
	re.Write(&reName)
	imName := "prefftimag.png"
	im.Write(&imName)
	
	inPlace := fftw.Array2{[...]int{fft.Rect.Dx(), fft.Rect.Dy()}, fft.Pix}
	
	fftw.FFT2To(&inPlace, &inPlace)
	
	return fft
}

func LogSpectrum(fft *FFTimage) (SippImage) {
	spect := new(SippGray)
	spect.Gray = image.NewGray(fft.Rect)
	spectPix := spect.Pix()
	temp := make ([]float64, len(spectPix))
	var max float64 = 0
	for index, pix := range fft.Pix {
		r := real(pix)
		i := imag(pix)
		val := math.Log(1 + math.Sqrt(r*r + i*i))
		if val > max {
			max = val
		}
		temp[index] = val
	}
	scale := 255.0 / max
	for index, pix := range temp {
		spectPix[index] = uint8(pix*scale)
	}
	return spect
}