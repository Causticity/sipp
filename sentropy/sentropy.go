// Copyright Raul Vera 2015-2016

// Package sentropy provides functions for computing entropy and delentropy
// of SippImages, as well as rendering these as images.
package sentropy

import (
	"image"
	"math"
)

import (
	. "github.com/Causticity/sipp/shist"
	. "github.com/Causticity/sipp/simage"
)

// SippEntropy is a structure that holds a reference to an image, the 1D
// histogram of its values, and the entropy values derived from the histogram.
type SippEntropy struct {
	// A reference to the image.
	Im SippImage
	// A reference to the 1D histogram
	Hist []uint32
	// The entropy for each bin value that actually occurred.
	BinEntropy []float64
	// The largest entropy value of any bin.
	MaxBinEntropy float64
	// The entropy for the image, i.e. the sum of the entropies for all
	// the bins.
	Entropy float64
}

// Entropy returns a SippEntropy structure for the given image.
func Entropy(im SippImage) (ent *SippEntropy) {
	ent = new(SippEntropy)
	ent.Im = im
	ent.Hist = GreyHist(im)
	total := float64(im.Bounds().Dx() * im.Bounds().Dy())
	normHist := make([]float64, len(ent.Hist))
	var check float64
	for i, binVal := range ent.Hist {
		normHist[i] = float64(binVal) / total
		check = check + normHist[i]
	}
	//fmt.Println("Normalised histogram sums to ", check)
	ent.BinEntropy = make([]float64, len(ent.Hist))
	for j, p := range normHist {
		if p > 0 {
			ent.BinEntropy[j] = -1.0 * p * math.Log2(p)
			ent.Entropy = ent.Entropy + ent.BinEntropy[j]
			if ent.BinEntropy[j] > ent.MaxBinEntropy {
				ent.MaxBinEntropy = ent.BinEntropy[j]
			}
		}
	}
	return
}

// EntropyImage returns a greyscale image of the entropy for each pixel.
func (ent *SippEntropy) EntropyImage() SippImage {
	entIm := new(SippGray)
	entIm.Gray = image.NewGray(ent.Im.Bounds())
	entImPix := entIm.Pix()

	// scale the entropy from (0-maxEnt) to (0-255)
	is16 := false
	if ent.Im.Bpp() == 16 {
		is16 = true
	}
	scale := 255.0 / ent.MaxBinEntropy
	width := ent.Im.Bounds().Dx()
	imPix := ent.Im.Pix()
	for y := 0; y < ent.Im.Bounds().Dy(); y++ {
		for x := 0; x < width; x++ {
			index := ent.Im.PixOffset(x, y)
			var val uint16 = uint16(imPix[index])
			if is16 {
				val = val<<8 | uint16(imPix[index+1])
			}
			entImPix[y*width+x] = uint8(math.Floor(ent.BinEntropy[val] * scale))
		}
	}
	return entIm
}

// SippDelentropy is a structure that holds a reference to a gradient histogram
// and the delentropy values derived from it.
type SippDelentropy struct {
	// A reference to the histogram we are computing from
	hist *SippHist
	// The delentropy for each bin value that actually occurred.
	binDelentropy []float64
	// The largest delentropy value of any bin.
	maxBinDelentropy float64
	// The delentropy for the image, i.e. the sum of the delentropies for all
	// the bins.
	Delentropy float64
}

// Delentropy returns a SippDelentropy structure for the given SippHist.
func Delentropy(hist *SippHist) (dent *SippDelentropy) {
	// Store the entropy values corresponding to the bin counts that actually
	// occurred.
	dent = new(SippDelentropy)
	dent.hist = hist
	dent.binDelentropy = make([]float64, hist.Max+1)
	total := float64(len(hist.Grad.Pix))
	for _, bin := range hist.Bin {
		if bin != 0 {
			// compute the entropy only once for a given bin value.
			if dent.binDelentropy[bin] == 0.0 {
				p := float64(bin) / total
				dent.binDelentropy[bin] = -1.0 * p * math.Log2(p)
			}
			dent.Delentropy += dent.binDelentropy[bin]
			if dent.binDelentropy[bin] > dent.maxBinDelentropy {
				dent.maxBinDelentropy = dent.binDelentropy[bin]
			}
		}
	}
	return
}

// HistDelentropyImage returns a greyscale image of the delentropy for each
// histogram bin.
func (dent *SippDelentropy) HistDelentropyImage() SippImage {
	// Make a greyscale image of the entropy for each bin.
	width, height := dent.hist.Size()
	dentGray := new(SippGray)
	dentGray.Gray = image.NewGray(image.Rect(0, 0, int(width), int(height)))
	dentGrayPix := dentGray.Pix()
	// scale the delentropy from (0-hist.maxBinDelentropy) to (0-255)
	scale := 255.0 / dent.maxBinDelentropy
	for i, val := range dent.hist.Bin {
		dentGrayPix[i] = uint8(dent.binDelentropy[val] * scale)
	}
	return dentGray
}

// DelEntropyImage returns a greyscale image of the entropy for each gradient
// pixel.
func (dent *SippDelentropy) DelEntropyImage() SippImage {
	// Make a greyscale image of the entropy for each bin.
	dentGray := new(SippGray)
	dentGray.Gray = image.NewGray(dent.hist.Grad.Rect)
	dentGrayPix := dentGray.Pix()
	// scale the entropy from (0-hist.maxBinDelentropy) to (0-255)
	scale := 255.0 / dent.maxBinDelentropy
	for i := range dentGrayPix {
		// The following lookup works as follows:
		// i - the gradient (and delentropy) image-pixel index
		// hist.BinIndex[i] - the histogram bin for that pixel
		// hist.Bin[hist.BinIndex[i] - the value in that bin
		// dent.binDelentropy[...] The delentropy for that value
		// We scale that delentropy and convert to an 8-bit pixel
		dentGrayPix[i] = uint8(dent.binDelentropy[dent.hist.Bin[dent.hist.BinIndex[i]]] * scale)
	}
	return dentGray
}
