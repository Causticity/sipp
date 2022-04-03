// Copyright Raul Vera 2015-2021

package shist

import (
	. "github.com/Causticity/sipp/simage"
)

const greyHistSize8BPP = 256
const greyHistSize16BPP = 65536

// GreyHist computes a 1D histogram of the greyscale values in the image.
func GreyHist(im SippImage) (hist []uint32) {
	histSize := greyHistSize8BPP
	is16 := false
	if im.Bpp() == 16 {
		histSize = greyHistSize16BPP
		is16 = true
	}

	hist = make([]uint32, histSize)
	imPix := im.Pix()
	for y := 0; y < im.Bounds().Dy(); y++ {
		for x := 0; x < im.Bounds().Dx(); x++ {
			index := im.PixOffset(x, y)
			var val uint16 = uint16(imPix[index])
			if is16 {
				val = val<<8 | uint16(imPix[index+1])
			}
			hist[val]++
		}
	}
	return
}