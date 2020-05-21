// Copyright Raul Vera 2015-2016

// Package sfft provides functions for the sipp package to compute an FFT and
// IFFT, as well as a function for displaying a spectrum as a grey-scale image.

package sfft

import (
	"image"
	"math"

    "github.com/davidkleiven/gosfft/sfft"
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

    ft := sfft.NewFFT2(fft.Rect.Dy(), fft.Rect.Dx())
    ft.FFT(fft.Pix)

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
