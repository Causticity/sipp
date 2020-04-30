// Copyright Raul Vera 2015-2016

// Package sfft provides functions for the sipp package to compute an FFT and
// IFFT using go-fftw, as well as a function for displaying a spectrum as a
// grey-scale image.

// By depending on FFTW, which uses the GPL, this package is also under the GPL.
// FFTW is linked via the go-fftw import below. The home page for FFTW is at
// http://www.fftw.org/

// Note that only this package and any program or package that imports it are
// under the GPL. The rest of the sipp packages are under the overall sipp
// license found in the README, which is the MIT license.

package sfft

import (
	"image"
	"math"

	"github.com/runningwild/go-fftw/fftw"
)

import (
	. "github.com/Causticity/sipp/scomplex"
	. "github.com/Causticity/sipp/simage"
)

type FFTImage struct {
	ComplexImage
}

func FFT(src SippImage) (fft *FFTImage) {
	comp := ToShiftedComplex(src)
	fft = &FFTImage{*comp}

	inPlace := fftw.Array2{[...]int{fft.Rect.Dx(), fft.Rect.Dy()}, fft.Pix}

	fftw.FFT2To(&inPlace, &inPlace)

	return fft
}

func LogSpectrum(fft *FFTImage) SippImage {
	spect := new(SippGray)
	spect.Gray = image.NewGray(fft.Rect)
	spectPix := spect.Pix()
	temp := make([]float64, len(spectPix))
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
		spectPix[index] = uint8(pix * scale)
	}
	return spect
}
