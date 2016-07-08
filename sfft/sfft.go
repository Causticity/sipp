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
	ComplexImage
}

func FFT(src SippImage) (fft *FFTimage) {
	comp := ToShiftedComplex(src)
	fft = &FFTimage {*comp}

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
		val := math.Log(1 + math.Hypot(real(pix), imag(pix)))
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